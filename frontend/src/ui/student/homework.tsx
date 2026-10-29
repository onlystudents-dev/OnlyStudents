import "../homework/homework.css";
import type {HomeworkSummary, Me} from "../../types/api.ts";
import {useState} from "preact/hooks";
import {useEffect} from "preact/compat";
import {fetchWithSchool} from "../../util/api.ts";
import Homework, {type PopupProps} from "../homework/homework.tsx";
import Overlay from "../../util/overlay/overlay.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faXmark} from "@fortawesome/free-solid-svg-icons";
import Attachment from "../attachment.tsx";
import { fromResponse, getKey } from "../../util/language.ts";
import SubmissionAttachmentInput from "../submission.tsx";
import { toast } from "react-toastify";

export default function StudentHomework({ me }: { me: Me }) {
    const [homework, setHomework] = useState<HomeworkSummary[] | null>(null)

    async function Fetch() {
        return fetchWithSchool("/api/v1/me/student/homework").then(async res => setHomework(await res.json() as HomeworkSummary[]))
    }

    useEffect(() => {
        void Fetch()
    }, [])

    return <Homework me={me} homework={homework} Popup={Popup} reloadSchoolStuff={Fetch} />
}

async function submitHomework(note: string, file_uuid: string, homework: HomeworkSummary) {
  await fetchWithSchool("/api/v1/me/student/homework/submit", {
      method: "POST", body: JSON.stringify({
          "file_uuid": file_uuid,
          "homework_id": homework.id,
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

function Popup({ homework: hw, clear }: PopupProps) {
    const homework = hw as HomeworkSummary

    return (
        <>
            <Overlay onClick={clear} />
            <div className="popup">
                <div className="popup-content">
                    <div className="absolute top-6 right-6 action-button hover:text-(--wrong-color)" onClick={clear}>
                        <FontAwesomeIcon icon={faXmark} />
                    </div>
                    <div className="flex flex-1 flex-col gap-3 min-h-0 overflow-y-auto p-6 pr-14">
                        <div className="text-3xl rubik truncate">{homework.title}</div>
                        <div className="text-2xl opacity-70 rubik truncate">{homework.subject}</div>
                        <div className="fredoka whitespace-pre-wrap wrap-break-word">{homework.description}</div>
                        <Attachment label={getKey("SUBMISSION")} missingLabel={getKey("NO_SUBMISSION")} attachment={{
                            note: homework.submission_note,
                            file_type: homework.submission_file_type,
                            file_size: homework.submission_file_size,
                            file_name: homework.submission_file_name,
                            file_uuid: homework.submission_file_uuid
                        }} />
                    </div>
                    <div className="shrink-0 px-6 pb-6 pr-14">
                        <SubmissionAttachmentInput notePlaceholder={getKey("HOMEWORK_SUBMISSION_NOTE_EXAMPLE")} label={getKey("HOMEWORK_SUBMISSION_SUBMIT")} onSubmit={(note: string, file_uuid: string) => {
                          void submitHomework(note, file_uuid, homework)
                        }}/>
                    </div>
                </div>
            </div>
        </>
    )
}
