import { useEffect, useState, type Dispatch, type StateUpdater } from "preact/hooks";
import { toast } from "react-toastify";
import { fromResponse, getKey } from "../util/language";
import { getFilesConfig, type FilesConfig } from "../auth/config";
import { fileCategory, KiB } from "../util/files";

type attachmentFile = {
  file_name: string,
  file_type: string,
  file_uuid: string,
  file_size: number,
}

type FileUploadResponse = {
    file_uuid: string,
}

async function FileUpload(setFile: Dispatch<StateUpdater<attachmentFile>>, setFileUploaded: Dispatch<StateUpdater<boolean>>, file: File, config: FilesConfig) {
    const match = config.allowed_exts.some(ext => file.name.endsWith(ext))

    if (!match) {
      toast.error(getKey("FILE_TYPE_NOT_ALLOWED"))
      return
    }

    if ((file.size / KiB) > config.max_file_size) {
      toast.error(getKey("FILE_TOO_BIG"))
      return
    }

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
      };
    });
}
export default function SubmissionAttachmentInput({ notePlaceholder, label, onSubmit }: { notePlaceholder: string, label?: string, onSubmit: (note: string, file_uuid: string) => void } ) {
    const [file, setFile] = useState<attachmentFile>({file_name: "", file_type: "", file_uuid: "", file_size: 0})
    const [fileUploaded, setFileUploaded] = useState(false)
    const [note, setNote] = useState("")
    const [previewUrl, setPreviewUrl] = useState<string | null>(null)
    const [config, setConfig] = useState<FilesConfig>()

    useEffect(() => {
        void getFilesConfig().then(setConfig)
    }, [setConfig])

    return (
        <>
            <div className="flex flex-col gap-3 border-t border-(--border-color) pt-3">
                <div className="text-xs font-medium uppercase tracking-wider opacity-50">{label}</div>
                <input type="text" placeholder={notePlaceholder} onChange={(e) => setNote(e.currentTarget.value)} />
                <input type="file" onChange={(e) => {
                  if (e.currentTarget.files) {
                      const file = e.currentTarget.files[0];
                      if (file) {
                          if (previewUrl) URL.revokeObjectURL(previewUrl)
                          setPreviewUrl(URL.createObjectURL(file))

                          if (!config) return

                          void FileUpload(setFile, setFileUploaded, file, config)
                      }
                  }
                }
                }></input>

                <button type="button" className="rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium hover:bg-white/10" onClick={() => {
                    if (!fileUploaded && note == "") return
                    onSubmit(note, file.file_uuid)
                }}>
                    {getKey("SUBMIT")}
                </button>
            </div>
        </>
    )
}
