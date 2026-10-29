import "./absences.css";
import type {Absence} from "../../types/api.ts";
import {useState} from "preact/hooks";
import {fetchInto, fetchWithSchool} from "../../util/api.tsx";
import {useCallback, useEffect} from "preact/compat";
import Skeleton from "../../util/skeleton/skeleton.tsx";
import {useReloadSchoolStuff} from "../../util/school.ts";
import Overlay from "../../util/overlay/overlay.tsx";
import {formatUnixDate} from "../../util/time.ts";
import {fromResponse, getKey, getLanguage} from "../../util/language.ts";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faXmark} from "@fortawesome/free-solid-svg-icons";
import Attachment from "../attachment.tsx";
import SubmissionAttachmentInput from "../submission.tsx";
import { toast } from "react-toastify";
import type {FeatureProps, FetchMe} from "../../types/props.ts";

async function submitAbsenceReason(note: string, fileUUID: string, absence: Absence, fetchMe: FetchMe): Promise<void> {
  await fetchWithSchool(fetchMe, "/api/v1/me/student/absence/submit_reason", {
    method: "POST", body: JSON.stringify({
      "file_uuid": fileUUID,
      "absence_id": absence.id,
      "note": note,
    }),
    headers: {
        "Content-Type": "application/json",
    },
  }).then(async response => {
      switch (response.status) {
          case 200: {
              location.reload()
              break
          }

          case 401: break
          default:
              toast.error(await fromResponse(response))
      }
  })
}

export default function StudentAbsences({ me, fetchMe }: FeatureProps) {
    const [absences, setAbsences] = useState<Absence[] | null>(null)
    const [current, setCurrent] = useState<Absence | null>(null)

    const Fetch = useCallback(async () => {
        const list = await fetchInto<Absence[]>(fetchMe, "/api/v1/me/student/absences")
        if (list) setAbsences(list)
    }, [fetchMe])

    useEffect(() => void Fetch(), [Fetch])
    useReloadSchoolStuff(() => void Fetch())

    return (
        <>
            <div className="w-full min-h-full flex flex-col items-center justify-start pt-8 gap-4">
                {current && <>
                    <Overlay onClick={() => setCurrent(null)} />
                    <div className={`popup ${!current.justified && "wrong"}`}>
                        <div className="popup-content">
                            <div className="absolute top-6 right-6 action-button hover:text-(--wrong-color)" onClick={() => setCurrent(null)}>
                                <FontAwesomeIcon icon={faXmark} />
                            </div>
                            <div className="flex flex-1 flex-col gap-3 min-h-0 overflow-y-auto p-6 pr-14">
                                <div className="flex flex-row items-baseline gap-3">
                                    <div className="text-3xl rubik flex-1 truncate">
                                        {current.subject}
                                    </div>
                                    <div className="shrink min-w-0 max-w-[45%]">
                                        {formatUnixDate(current.date, getLanguage(me)?.key || "en-US")}
                                    </div>
                                </div>
                                <div className="text-2xl opacity-70 rubik truncate">
                                    {current.teacher_first_name != "" ? `${current.teacher_first_name} ${current.teacher_last_name}` : current.teacher_last_name}
                                </div>
                                <Attachment label={getKey("GIVEN_REASON")} missingLabel={getKey("NO_GIVEN_REASON")} attachment={{ ...current, note: current.reason_note }} />
                            </div>
                            <div className="shrink-0 px-6 pb-6 pr-14">
                              <SubmissionAttachmentInput notePlaceholder={getKey("ABSENCE_REASON_NOTE_EXAMPLE")} label={getKey("ABSENCE_JUSTIFICATION_SUBMIT")} onSubmit={(note: string, fileUUID: string) => void submitAbsenceReason(note, fileUUID, current, fetchMe)} />
                            </div>
                        </div>
                    </div>
                </>}
                {absences?.map((absence) => (
                    <div key={absence.id} className={`absence ${!absence.justified && "wrong"}`} onClick={() => setCurrent(absence)}>
                        <div className="flex flex-col items-start justify-between">
                          <p className="max-w-100 truncate">
                              {absence.subject}
                          </p>
                          <p className="max-w-100 truncate opacity-70">
                              {absence.reason_note != "" ? absence.reason_note :
                               absence.file_name != "" ? absence.file_name :
                               `${absence.teacher_first_name} ${absence.teacher_last_name} (${formatUnixDate(absence.date, getLanguage(me)?.key || "en-US")})`}
                          </p>

                        </div>
                        <div className="flex flex-col items-start justify-between">

                        </div>
                    </div>
                )) || Array.from({ length: 9 }).map((_, i) => (
                    <Skeleton width={640} height={64} color="var(--hover-color)" key={i} className="border-b-2 border-r-2 border-(--border-color) rounded-3xl!" />
                ))}
            </div>
        </>
    )
}
