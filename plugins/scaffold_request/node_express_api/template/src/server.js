import express from 'express'

const app = express()
const port = Number(process.env.PORT || [[PORT]])

app.get('/healthz', (_req, res) => {
  res.json({ status: 'ok', service: '[[SERVICE_NAME]]' })
})

app.listen(port, () => {
  console.log(`[[SERVICE_NAME]] listening on ${port}`)
})
