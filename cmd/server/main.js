const ws = require('nodejs-websocket')

const connList = [null, null]
let index = 0

const createServer = () => {
  let server = ws.createServer(connection => {
    connList[index] = connection
    connection.myID = index
    if (index == 0) {
      index = 1
    } else {
      index = 0
    }
    console.log('income')
    connection.on('text', function(result) {
      try {
        // for (let i = 2; i <= 5; i++) {
        //   const input = JSON.parse(result)
        //   input.clientID = i
        //   const text = JSON.stringify(input)
        //   connection.send(text)
        //   console.log('发送消息', connection.myID, text)
        // }
        const input = JSON.parse(result)
        const text = JSON.stringify(input)
        const otherIndex = connection.myID === 0 ? 1 : 0
        const other = connList[otherIndex]
        if (other) {
          other.send(text)
        }
      } catch (err) {
        console.error(err)
      }
      console.log('收到消息', connection.myID, result)
    })
    connection.on('connect', function(code) {
      console.log('开启连接', code)
    })
    connection.on('close', function(code) {
      console.log('关闭连接', code)
    })
    connection.on('error', function(code) {
      console.log('异常关闭', code)
    })
  })
  return server
}

const server = createServer()
server.listen(8966)
