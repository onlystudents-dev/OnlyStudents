export const KiB = 1 << 10
const MiB = 1 << 20
const GiB = 1 << 30
const max = 1 << 10
export function fileCategory(file: File): string {
  if (file.type.startsWith("image/")) return "image"
  if (file.type.includes("presentation") || /\.(ppt|pptx)$/i.test(file.name)) return "presentation"
  if (file.type.includes("spreadsheet") || /\.(xls|xlsx|csv)$/i.test(file.name)) return "spreadsheet"
  if (/\.(zip|rar|7z|tar|gz)$/i.test(file.name)) return "archive"
  return "document"
}

export function formatFileSize(size: number): string {
  if (size > (MiB * max)) {
    return `${(size / GiB).toFixed(1)} GiB`
  }

  if (size > (KiB * max)) {
    return `${(size / MiB).toFixed(1)} MiB`
  }

  return `${(size / KiB).toFixed(1)} KiB`
}
