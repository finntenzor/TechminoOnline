const ws = require('nodejs-websocket')

const connList = {}
let nextID = 1
let nextGarbageID = 1

const send = (client, msg) => {
  client.send(msg)
  console.log('发送消息', client.clientID, msg)
}

const sendToAll = msg => {
  for (let i = 1; i <= 5; i++) {
    const other = connList[i]
    if (other) {
      send(other, msg)
    }
  }
}

const sendExcept = (msg, except) => {
  for (let i = 1; i <= 5; i++) {
    const other = connList[i]
    if (other && other.clientID !== except) {
      send(other, msg)
    }
  }
}

const randTarget = (except) => {
  let target = 0
  for (let i = 0; i < 1000; i++) {
    target = Math.ceil(Math.random() * (nextID - 1))
    if (target !== except) {
      return target
    }
  }
  return (except + 1) % (nextID - 1) + 1
}

const createServer = () => {
  let server = ws.createServer(connection => {
    connList[nextID] = connection
    connection.clientID = nextID
    console.log('income', nextID)
    nextID++

    // 连接的时候 告诉客户端它的ID
    connection.send(JSON.stringify({
      messageType: 1,
      clientID: connection.clientID
    }))
    // 收到某个客户端发来的包
    connection.on('text', function(result) {
      try {
        console.log('收到消息', connection.clientID, result)
        const input = JSON.parse(result)
        const message = {
          ...input,
          clientID: connection.clientID
        }
        if (input.messageType > 1 && input.messageType < 8) {
          // 移动动作
          // 给除了自己以外的人发移动动作
          sendExcept(JSON.stringify(message), connection.clientID)
        } else if (input.messageType === 8) {
          // 垃圾行包1
          message.garbageID = nextGarbageID
          nextGarbageID++
          // message.from = connection.clientID // 似乎是不必要的
          message.target = randTarget(message.clientID) // 随机选一个已上线目标 现在客户端有问题 之后完全由客户端选目标
          message.position = Math.ceil(Math.random() * 10)
          // 给所有人发 垃圾行包1
          sendToAll(JSON.stringify(message))
        } else if (input.messageType === 9) {
          // 垃圾行包2

          // 因为客户端是先涨缓冲区 再发包
          // 所以给这个人以外的所有人发包
          sendExcept(JSON.stringify(message), connection.clientID)
        } else if (input.messageType === 10) {
          // 垃圾行包3

          // 同理 告诉其他人要涨垃圾行 也是给自己以外的人发
          sendExcept(JSON.stringify(message), connection.clientID)
        }
      } catch (err) {
        console.error(err)
      }
    })
    connection.on('connect', function(code) {
      console.log('开启连接', connection.clientID, code)
    })
    connection.on('close', function(code) {
      console.log('关闭连接', connection.clientID, code)
      // 从map中删除连接
      delete connList[connection.clientID]
      // 全部关闭的时候id重置为1开始
      if (Object.keys(connList).length === 0) {
        console.log('状态重置')
        nextID = 1
        nextGarbageID = 1
      }
    })
    connection.on('error', function(err) {
      console.log('异常关闭', connection.clientID, err.message)
    })
  })
  return server
}

const server = createServer()
server.listen(8966)
