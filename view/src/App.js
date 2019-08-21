import React from 'react';
import './App.css';
import { List, Comment, Input, Card, Row, Col, Button, message, notification } from 'antd';

class CommonMsg extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      name: "",
      data: [],
      msg: "",
    };
  }

  socket;
  componentDidMount() {
    if (window["WebSocket"]) {
      const loading = message.loading("连接服务器中。。", 0);
      this.socket = new WebSocket("wss://lwlgo.com:12345/ws");
      this.socket.onopen = function (evt) {
        loading();
      }
      this.socket.onerror = function (evt) {
        loading();
        setTimeout(() => { message.error("Connection Failed.") }, 500);
      }
      const addMsg = (msg) => {
        this.setState({
          data: this.state.data.concat(msg),
        });
        const contentObj = document.getElementById('content');
        contentObj.scrollTop = contentObj.scrollHeight;
      };
      this.socket.onclose = function (evt) {
        loading();
        setTimeout(() => { message.info("Connection Closed.") }, 600);
      }
      this.socket.onmessage = function (evt) {
        const playload = JSON.parse(evt.data)
        switch (playload.type) {
          case "broadcast":
            const datas = [
              {
                author: playload.name,
                content: playload.content,
              },
            ];
            addMsg(datas);
            break;
          case "success":
            message.success(playload.content)
            break;
          case "error":
            message.error(playload.content)
            break;
          case "info":
            message.info(playload.content)
            break;
          case "notification":
            notification.open({
              message: '通知',
              description:
                playload.content,
            });
            break;
          default:
            message.info("接收到未知消息")
            break;
        }
      }
    } else {
      message.error("WebSockets Not Support.");
    }

  }
  componentWillUnmount() {
    if (this.socket) {
      this.socket.close();
    }
  }
  sendMsg(msg, no) {
    if (!this.socket) {
      return false;
    }
    var data = msg;
    if (!data) {
      return false;
    }
    let dataLen = data.length * 4 + 2;
    var bf = new ArrayBuffer(dataLen);
    if (bf.byteLength !== dataLen) {
      console.log(bf.byteLength + "!=" + dataLen);
      setTimeout(() => { message.error("申请内存失败") }, 500);
      return false;
    }
    var dataView = new DataView(bf);
    dataView.setUint16(0, no);//协议号 1 设置名称
    if (data.length > 0) {
      for (let i = 0; i < data.length; i++) {
        dataView.setInt32(i * 4 + 2, data.charCodeAt(i));
      }

    }
    this.socket.send(bf);
    this.setState({
      name:"",
      msg:"",
    });
    return true;
  }
  handleNameChange(e) {
    this.setState(
      {
        name: e.target.value,
      }
    );
  }
  handleMsgChange(e) {
    this.setState(
      {
        msg: e.target.value,
      }
    );
  }
  render() {
    return (
      <div style={{ background: '#ECECEC', padding: '30px', height: "100%", width: "100%", position: "absolute" }}>
        <Row style={{ height: "90%" }}>
          <Card bordered={false} id="content" style={{ height: "100%", overflowX: "hidden", overflowY: "scroll", }} loading={false}>
            <List

              style={{ textAlign: "left", }}
              className="comment-list"
              itemLayout="horizontal"
              dataSource={this.state.data}
              renderItem={
                item => (
                  <li>
                    <Comment
                      author={item.author}
                      content={item.content}
                    />
                  </li>
                )
              }
            >
            </List>

          </Card>

        </Row>
        <Row style={{ paddingTop: "10px" }}>
          <Col span={2}><Input placeholder="设置名称" value={this.state.name} onChange={this.handleNameChange.bind(this)}></Input></Col>
          <Col span={1}><Button type="primary" onClick={() => { this.sendMsg(this.state.name, 1) }}>设置</Button></Col>
        </Row>
        <Row style={{ paddingTop: "10px" }}>
          <Col span={8}><Input placeholder="需要发送的消息" value={this.state.msg} onChange={this.handleMsgChange.bind(this)}></Input></Col>
          <Col span={1}><Button type="primary" onClick={() => { this.sendMsg(this.state.msg, 2) }}>发送</Button></Col>
        </Row>
      </div>
    );

  }
}

function App() {
  return (
    <div className="App" >
      <CommonMsg />
    </div>
  );
}

export default App;
