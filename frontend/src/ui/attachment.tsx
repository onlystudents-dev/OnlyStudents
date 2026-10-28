import {
    faFileLines,
    faFilePowerpoint,
    faFileImage,
    faDownload,
    faFile,
    faFileExcel,
    faFileArchive,
    faFileAudio,
} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {formatFileSize} from "../util/files";

type attachmentType = {
  note: string,
  file_type: string,
  file_size: number
  file_name: string,
  file_uuid: string,
}

export default function Attachment({ attachment, label, missingLabel }: { attachment: attachmentType, label: string, missingLabel: string } ) {
  const hasNote = attachment.note != ""
  const hasFile = attachment.file_name != ""

  return (
      <div className="flex flex-col gap-3 border-t border-(--border-color) pt-3">
        <div className="text-xs font-medium uppercase tracking-wider opacity-50">{label}</div>
        {(!hasNote && !hasFile) && (
          <div className="text-xs font-medium uppercase tracking-wider opacity-50">{missingLabel}</div>
        )}
        {hasNote && (
            <div className="fredoka whitespace-pre-wrap wrap-break-word">
                {attachment.note}
            </div>
        )}
        {hasFile && (
            <div className="flex flex-row items-center gap-3 bg-white/5 border border-white/10 rounded-lg p-3">
                <div className={`flex items-center justify-center w-12 h-12 rounded-lg text-2xl ${
                    attachment.file_type === "document" ? "bg-blue-500/20 text-blue-400" :
                    attachment.file_type === "presentation" ? "bg-orange-500/20 text-orange-400" :
                    attachment.file_type === "image" ? "bg-green-500/20 text-green-400" :
                    attachment.file_type === "spreadsheet" ? "bg-purple-500/20 text-purple-400" :
                    "bg-gray-500/20 text-gray-400"
                }`}>
                    <FontAwesomeIcon
                        icon={
                            attachment.file_type === "document" ? faFileLines :
                            attachment.file_type === "presentation" ? faFilePowerpoint :
                            attachment.file_type === "image" ? faFileImage :
                            attachment.file_type === "spreadsheet" ? faFileExcel :
                            attachment.file_type === "audio" ? faFileAudio :
                            attachment.file_type === "archive" ? faFileArchive :
                            attachment.file_type === "image" ? faFileImage :
                            faFile
                        }
                    />
                </div>
                <div className="flex-1 min-w-0">
                    <div className="fredoka text-lg truncate">{attachment.file_name}</div>
                    <div className="text-sm opacity-60">{formatFileSize(attachment.file_size)}</div>
                </div>
                <a href={`/uploads/${attachment.file_uuid}`} className="action-button w-10 h-10 flex items-center justify-center rounded-lg bg-white/10 hover:bg-white/20 transition" download>
                    <FontAwesomeIcon icon={faDownload} />
                </a>
            </div>
        )}
      </div>
  )
}
