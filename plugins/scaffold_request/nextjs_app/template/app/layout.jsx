export const metadata = {
  title: '[[SERVICE_NAME]]',
}

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}
