export const makeInvisibleDownloadAnchor = (data: Blob, filename: string) => {
  const link = document.createElement("a")
  link.href = window.URL.createObjectURL(data)
  link.setAttribute("download", filename)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}
