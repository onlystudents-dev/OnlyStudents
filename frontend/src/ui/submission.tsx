import { useState, type Dispatch, type StateUpdater } from "preact/hooks";
import { toast } from "react-toastify";
import { fromResponse } from "../util/language";

type attachmentFile = {
  file_name: string
  file_type: string
  file_uuid: string
  file_size: number
}

type FileUploadResponse = {file_uuid: string}

function fileCategory(file: File): string {
  if (file.type.startsWith("image/")) return "image"
  if (file.type.includes("presentation") || /\.(ppt|pptx)$/i.test(file.name)) return "presentation"
  if (file.type.includes("spreadsheet") || /\.(xls|xlsx|csv)$/i.test(file.name)) return "spreadsheet"
  if (/\.(zip|rar|7z|tar|gz)$/i.test(file.name)) return "archive"
  return "document"
}

async function FileUpload(setFile: Dispatch<StateUpdater<attachmentFile>>, setFileUploaded: Dispatch<StateUpdater<boolean>>, file: File) {
  const data = new FormData()
  data.append('file', file)

  await fetch("/api/v1/files", { method: "POST", body: data }).then(async response => {
      switch (response.status) {
          case 200: {
              const json = await response.json() as FileUploadResponse;

              setFile({
                  file_name: file.name,
                  file_type: fileCategory(file),
                  file_uuid: json.file_uuid,
                  file_size: file.size,
              })
              setFileUploaded(true)

              break
          }

          case 401: break
          default:
              toast.error(await fromResponse(response))
      }
  })
}

export default function SubmissionAttachmentInput({ label, onSubmit }: { label?: string, onSubmit: (note: string, file_uuid: string) => void } ) {
  const [file, setFile] = useState<attachmentFile>({file_name: "", file_type: "", file_uuid: "", file_size: 0})
  const [fileUploaded, setFileUploaded] = useState(false)
  const [note, setNote] = useState("")
  const [previewUrl, setPreviewUrl] = useState<string | null>(null)

  return (
    <>
      <div className="flex flex-col gap-3 border-t border-(--border-color) pt-3">
        {label && (
            <div className="text-xs font-medium uppercase tracking-wider opacity-50">{label} Submit</div>
        )}
      <input type="text" placeholder="Your justification here..." onChange={(e) => setNote(e.currentTarget.value)} />
      <input type="file" onChange={(e) => {
        if (e.currentTarget.files) {
          const file = e.currentTarget.files[0];
          if (file) {
            if (previewUrl) URL.revokeObjectURL(previewUrl)
            setPreviewUrl(URL.createObjectURL(file))
            void FileUpload(setFile, setFileUploaded, file);
          };
        };
      }
      }></input>

      <button type="button" className="rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium hover:bg-white/10" onClick={() => {
        if (!fileUploaded && note == "") {
          return
        }
        onSubmit(
          note,
          file.file_uuid
        )
      }}>
          Submit
      </button>
      </div>
    </>
  )
}
