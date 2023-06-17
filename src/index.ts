import express from 'express';

const app = express()

app.get('/foo', (req, res) => {
    res.send({
        ok: true,
        foo: 'bar'
    })
})

app.listen(9090, () => {
    console.log('server running on port 9090')
})