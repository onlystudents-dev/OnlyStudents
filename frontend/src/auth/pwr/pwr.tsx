import "./pwr.css";
import Button from "../../util/button/button.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faArrowLeft, faPaperPlane, faPlane} from "@fortawesome/free-solid-svg-icons";
import React, {useCallback, useEffect, useState} from "preact/compat";
import {toast} from "react-toastify";
import PasswordCheck from "../pwc.tsx";
import {fromResponse, getKey} from "../../util/language.ts";
import {DEFAULT_CONFIG, getPasswordConfig, type PasswordConfig} from "../config.ts";
import RoleSelect from "../roleSelect.tsx";
import {useCountdown} from "../../util/countdown.ts";
import {getRetryAfter, postJSON} from "../../util/api.ts";
import type { ForgetPasswordCheckResponse, ForgetPasswordConfirmResponse } from "../../types/api.ts";
import { MnemonicSecrettoSeed, normalizeMnemonicWords, validateMnemonicSecret } from "../../util/crypto/mnemonic.ts";
import { b64decode, b64encode, constantTimeEqual, decodeOpaqueExportKey, deriveWrapKey, parseParams, recoveryDigest, unwrapMK, wrapMK, type Bytes } from "../../util/crypto/keys.ts";
import { opaqueRegister } from "../opaque.ts";

export default function PasswordReset({ resetActive, onClose, role, setRole, id, setId, red, checkUserID, setWaiting, setResetSent }: {resetActive: boolean, onClose: () => void, role: string, setRole: React.Dispatch<React.SetStateAction<string>>, id: string, setId: React.Dispatch<React.SetStateAction<string>>, red: boolean, checkUserID: (id: string) => void, setWaiting: React.Dispatch<React.SetStateAction<boolean>>, setResetSent:  React.Dispatch<React.SetStateAction<boolean>>}) {
    const [codeSent, setCodeSent] = useState(false)
    const [transitioning, setTransitioning] = useState(false)

    const [code, setCode] = useState("")

    const [recoveryKey, setRecoveryKey] = useState("")
    const [recoveryKeyWrong, setRecoveryKeyWrong] = useState(false)

    const [password, setPassword] = useState("")
    const [confirmPassword, setConfirmPassword] = useState("")

    const [passed, setPassed] = useState(false)

    const [config, setConfig] = useState<PasswordConfig | null>(null)

    const [remaining, updateRemaining] = useCountdown()

    useEffect(() => {
        void getPasswordConfig().then(setConfig)
    }, []);

    const submit = useCallback(async (url: string, body: unknown, onOk: (response: Response) => void) => {
        setWaiting(true)
        const response = await postJSON(url, body)

        if (response.status === 200) {
            onOk(response)
        } else {
            if (response.status === 429) updateRemaining(getRetryAfter(response))
            toast.error(await fromResponse(response))
        }
        setWaiting(false)
    }, [setWaiting, updateRemaining])

    const sendResetEmail = useCallback(async () => {
        if (red || remaining || !id) return
        await submit("/api/v1/forget_password", { user: Number(id), role }, () => {
            setTransitioning(true)
            setResetSent(true)
            setTimeout(() => setCodeSent(true), 200)
        })
    }, [red, remaining, id, role, setResetSent, submit])

    const resetPassword = useCallback(async () => {
      if (Number(remaining) > 0 || !code || !password || !passed || password !== confirmPassword) return

      const recoveryKeyNormalized = normalizeMnemonicWords(recoveryKey)

      if (!validateMnemonicSecret(recoveryKeyNormalized)) {
        setRecoveryKeyWrong(true)
        return
      }

      setRecoveryKeyWrong(false)

      await submit("/api/v1/forget_password_check", { reset_code: code }, (checkResponse) => {
        void (async () => {
          let unwrappedMk: Bytes
          try {
            const json = await checkResponse.json() as ForgetPasswordCheckResponse

            const params = parseParams(json.recovery_key_params)
            const digest = recoveryDigest(recoveryKeyNormalized, b64decode(json.recovery_key_salt), {
              ...DEFAULT_CONFIG.e2ee,
              argon2_time: params.t,
              argon2_memory: params.m,
              argon2_threads: params.p,
              argon2_key_len: params.dkLen,
            })

            if (!constantTimeEqual(digest, b64decode(json.recovery_key_hash))) {
              setRecoveryKeyWrong(true)
              return
            }

            unwrappedMk = await unwrapMK(
              b64decode(json.wrapped_mk_recovery),
              await deriveWrapKey(MnemonicSecrettoSeed(recoveryKeyNormalized), "recovery"),
            )
          } catch {
            setRecoveryKeyWrong(true)
            return
          }

          await submit("/api/v1/forget_password_confirm", { reset_code: code }, (confirmResponse) => {
            void (async () => {
              let json: ForgetPasswordConfirmResponse
              let res: Awaited<ReturnType<typeof opaqueRegister>>
              try {
                json = await confirmResponse.json() as ForgetPasswordConfirmResponse
                res = await opaqueRegister({ enrollToken: json.enroll_token, password })
              } catch {
                toast.error(getKey("RECOVERY_FAILED"))
                return
              }

              if (res.status === "ok") {
                try {
                  const wrapped = await wrapMK(
                    unwrappedMk,
                    await deriveWrapKey(decodeOpaqueExportKey(res.exportKey), "password"),
                  )
                  await submit(
                    "/api/v1/me/e2ee/rewrap",
                    { crypto_grant_handle: json.crypto_grant_handle, wrapped_mk_password: b64encode(wrapped) },
                    () => onClose(),
                  )
                } catch {
                  toast.error(getKey("RECOVERY_FAILED"))
                }
                return
              }

              if (res.status === "ratelimited") {
                const retryAfter = getRetryAfter(confirmResponse)
                toast.error(getKey("TOO_MANY_REQUESTS", String(retryAfter)))
                updateRemaining(retryAfter)
                return
              }

              toast.error(res.error ? getKey(res.error) : getKey("RECOVERY_FAILED"))
            })()
          })
        })()
      })
    }, [remaining, code, password, passed, confirmPassword, recoveryKey,
        setRecoveryKeyWrong, updateRemaining, onClose, submit])

    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === "Enter") void (codeSent ? resetPassword() : sendResetEmail())
        }

        window.addEventListener("keydown", handleKeyDown)

        return () => {
            window.removeEventListener("keydown", handleKeyDown)
        }
    }, [codeSent, sendResetEmail, resetPassword])

    return (
        <>
            <div className={`content in outback ${!resetActive && "hid"} ${transitioning ? "h-170" : "h-86"}`}>
                <h1 className="self-center text-5xl font-bold mb-8 rubik">{getKey("FORGOT_PASSWORD_TITLE")}</h1>
                {!codeSent && <div className={`loginput fredoka out ${transitioning && "hid"}`}>
                    <RoleSelect value={role} onChange={setRole} />
                    <input className={`${red && "wrong"}`} type="text" placeholder="User ID" value={id} onChange={(e) => {setId(e.currentTarget.value); checkUserID(e.currentTarget.value)}} />
                </div>}
                {codeSent && <div className="loginput fredoka in">
                    <input type="text" placeholder={getKey("CODE")} onChange={(e) => setCode(e.currentTarget.value)} />
                    <input type="text" className={`${recoveryKeyWrong && "wrong"}`} placeholder={getKey("RECOVERY_KEY")} onChange={(e) => setRecoveryKey(e.currentTarget.value)} />
                    <input type="password" placeholder={getKey("PASSWORD")} maxLength={config?.max_length} onChange={(e) => setPassword(e.currentTarget.value)} />
                    <input className={`${password !== confirmPassword && "wrong"}`} type="password" placeholder={getKey("CONFIRM_PASSWORD")} maxLength={config?.max_length} onChange={(e) => setConfirmPassword(e.currentTarget.value)} />
                    <PasswordCheck password={password} confirmPassword={confirmPassword} setPassed={setPassed} />
                </div>}

                <div className="pwrbutton">
                    <Button onClick={onClose}>
                        <FontAwesomeIcon icon={faArrowLeft} /> {getKey("BACK")}
                    </Button>
                    <div className="flex flex-row items-center gap-4">
                        <p className="text-(--wrong-color) text-3xl fredoka">{remaining}</p>
                        {!codeSent && <Button onClick={sendResetEmail} className={`outbottom ${transitioning && "hid"}`}>
                            {getKey("SEND_EMAIL")} <FontAwesomeIcon icon={faPaperPlane} />
                        </Button>}
                        {codeSent && <Button onClick={resetPassword} className="in">
                            {getKey("RESET")} <FontAwesomeIcon icon={faPlane} />
                        </Button>}
                    </div>
                </div>
            </div>
        </>
    )
}
