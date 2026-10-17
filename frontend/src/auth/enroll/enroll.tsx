import "./enroll.css";
import {fromResponse, getKey} from "../../util/language.ts";
import Button from "../../util/button/button.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faArrowRight} from "@fortawesome/free-solid-svg-icons";
import {useCallback, useEffect, useState} from "preact/compat";
import PasswordCheck from "../pwc.tsx";
import {getE2EEConfig, getPasswordConfig, type PasswordConfig} from "../authConfig.ts";
import {opaqueRegister} from "../opaque.ts";
import {toast} from "react-toastify";
import Loading from "../../util/loading.tsx";
import { decodeOpaqueExportKey, deriveWrapKey, generateMK, paramsToString, randomBytes, recoveryDigest, wrapMK } from "../../util/crypto/keys.ts";
import { generateMnemonicSecret, MnemonicSecrettoSeed, normalizeMnemonicSecret } from "../../util/crypto/mnemonic.ts";
import { getRetryAfter } from "../../util/api.ts";

const CHUNK = 0x8000;

export function b64encode(u8: Uint8Array): string {
  let s = "";
  for (let i = 0; i < u8.length; i += CHUNK) {
    s += String.fromCharCode(...u8.subarray(i, i + CHUNK))
  }
  return btoa(s)
}

export default function Enroll() {
    const enrollToken = new URLSearchParams(location.search).get("enroll_token");

    const [password, setPassword] = useState("")
    const [confirmPassword, setConfirmPassword] = useState("")

    const [passed, setPassed] = useState(false)
    const [waiting, setWaiting] = useState(false)
    const [remaining, setRemaining] = useState("")
    const [recovery, setRecovery] = useState(false)
    const [recoveryMnemonic, setRecoveryMnemonic] = useState("")
    const [accepted, setAccepted] = useState(false)
    const [revealed, setRevealed] = useState(false)

    const [config, setConfig] = useState<PasswordConfig | null>(null)

    async function PrepareE2EE(exportKey: string) {
      const e2ee_config = await getE2EEConfig()
      const MK = generateMK()
      const mnemonic = generateMnemonicSecret(e2ee_config.recovery_key_word_count)
      const OpaqueWrappedMK = await wrapMK(MK, await deriveWrapKey(decodeOpaqueExportKey(exportKey), "password"))
      const MnemonicWrappedMK = await wrapMK(MK, await deriveWrapKey(MnemonicSecrettoSeed(normalizeMnemonicSecret(mnemonic)), "recovery"))
      const RecoveryKeySalt = randomBytes(e2ee_config.argon2_salt_len)
      setWaiting(true)
      try {
          const response = await fetch("/api/v1/me/e2ee/enroll", {
              method: "POST",
              headers: {
                  "Content-Type": "application/json",
              },
              body: JSON.stringify({
                  wrapped_mk_password: b64encode(OpaqueWrappedMK),
                  wrapped_mk_recovery: b64encode(MnemonicWrappedMK),
                  recovery_key_hash: b64encode(recoveryDigest(mnemonic, RecoveryKeySalt, e2ee_config)),
                  recovery_key_salt: b64encode(RecoveryKeySalt),
                  recovery_key_params: paramsToString(e2ee_config)
              })
          })

          switch (response.status) {
              case 200:
                  setRecoveryMnemonic(mnemonic)
                  setRecovery(true)
                  break
              case 429: {
                toast.error(getKey("TOO_MANY_REQUESTS", String(getRetryAfter(response))))
                setRemaining(String(getRetryAfter(response)))
                break
              }
              default:
                toast.error(await fromResponse(response))
          }
      } finally { setWaiting(false)}
    }

    const register = useCallback(async () => {
        if (!passed || remaining || !enrollToken) return

        setWaiting(true)

        try {
            const res = await opaqueRegister({ enrollToken, password })

            switch (res.status) {
                case "ok":
                    await PrepareE2EE(res.exportKey)
                    break
                case "ratelimited":
                    toast.error(getKey("TOO_MANY_REQUESTS", String(res.retryAfter)))
                    setRemaining(String(res.retryAfter))
                    break
                case "error":
                    toast.error(
                        res.error
                            ? getKey(res.error)
                            : getKey("ENROLL_FAILED")
                    )
                    break
            }
        } finally {
            setWaiting(false)
        }
    }, [passed, remaining, enrollToken, password])

    useEffect(() => {
        void getPasswordConfig().then(setConfig)

        const handleKeyDown = (e: KeyboardEvent) => {
            if (recovery) {
                if (e.key === "Enter" || e.key === " ") setRevealed(true)
            } else {
                if (e.key === "Enter") void register()
            }
        }

        window.addEventListener("keydown", handleKeyDown)

        return () => {
            window.removeEventListener("keydown", handleKeyDown)
        }
    }, [recovery, register])

    return (
        <>
            <div className="cantar ebox">
                <div className="econtent">
                    <h1 className="self-center text-5xl font-bold mb-4 rubik">
                        {getKey("ENROLL")}
                    </h1>

                    {!recovery && <div className="loginput fredoka">
                      <input type="password" placeholder={getKey("PASSWORD")} maxLength={config?.max_length} value={password} onChange={(e) => setPassword(e.currentTarget.value)} />
                      <input className={`${password !== confirmPassword && "wrong"}`} type="password" placeholder={getKey("CONFIRM_PASSWORD")} maxLength={config?.max_length} value={confirmPassword} onChange={(e) => setConfirmPassword(e.currentTarget.value)} />
                      <PasswordCheck password={password} confirmPassword={confirmPassword} setPassed={setPassed} />
                    </div>}

                    {recovery && <div className="loginput fredoka">
                        <div>{getKey("RECOVERY_KEY_SHOW")}</div>

                        <div className={`recovery-key ${!revealed && "spoiler"}`} onClick={() => { if (revealed) void copy(); else setRevealed(true); }} role="button" tabIndex={0} title={!revealed && getKey("RECOVERY_KEY_SHOW") || undefined}>
                            {recoveryMnemonic}
                        </div>

                        <div>{getKey("RECOVERY_KEY_STORAGE")}</div>

                        {revealed && <label className="recovery-accept">
                            <input type="checkbox" checked={accepted} onChange={(e) => setAccepted(e.currentTarget.checked)} />
                            {getKey("RECOVERY_KEY_ACCEPT")}
                        </label>}
                    </div>}

                    <div className="elogbutton">
                        <p className="text-(--wrong-color) text-3xl fredoka">
                            {remaining}
                        </p>

                        {recovery ? <Button className={accepted ? "hover-shadow" : "opacity-50 cursor-not-allowed!"} onClick={() => { if (accepted) location.href = "/me" }}>
                            {getKey("ENROLL")} <FontAwesomeIcon icon={faArrowRight} />
                        </Button> : <Button className="hover-shadow" onClick={register}>
                            {getKey("CONTINUE")} <FontAwesomeIcon icon={faArrowRight} />
                        </Button>}
                    </div>
                </div>
            </div>

            {waiting && <Loading />}
        </>
    )

    async function copy() {
        await navigator.clipboard.writeText(recoveryMnemonic)
        toast.success(getKey("COPIED"))
    }
}
