import { h, render, Component } from 'preact'
import Fuse from 'fuse-js-latest'
import Modal from 'react-modal'
import ReactHintFactory from 'react-hint'
import axios from "axios"
import ws from 'ws'

const ReactHint = ReactHintFactory({createElement: h, Component})

const customStyles = {
  content : {
    top                   : '50%',
    left                  : '50%',
    right                 : 'auto',
    bottom                : 'auto',
    marginRight           : '-50%',
    transform             : 'translate(-50%, -50%)'
  }
}

Modal.setAppElement('#chat')

const months = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec'];

function formatDate(timestamp) {
  const d = new Date(timestamp*1000)
  const formated = {
    day: d.getDate() + ' ' + months[d.getMonth()],
    time: d.getHours() + ':' + d.getMinutes()
  }
  return formated
}

const axiosConfig = {
  headers: {
    'content-Type': 'application/json',
    'Accept': '/',
    'Cache-Control': 'no-cache',
  },
  credentials: 'same-origin'
}

const fuseOptions = {
  shouldSort: true,
  threshold: 0.6,
  location: 0,
  distance: 100,
  maxPatternLength: 12,
  minMatchCharLength: 1,
  keys: [
    "email",
    "name",
    "username",
  ]
}

axios.defaults.withCredentials = true

/** @jsx h */
class BroadcastModal extends Component {
  constructor(props) {
    super(props)
    this.state = {
      broadcasting: false,
      message: "",
    }
    this.canClose = this.canClose.bind(this)
    this.handleChange = this.handleChange.bind(this)
  }
  handleChange(event) {
    this.setState({text: event.target.value})
  }
  componentDidMount() {
    // FIXME: duplicate websocket connection.

  }
  canClose() {
    // TODO: render error msg and disable button if broadcasting
    if (this.broadcasting) return
    this.props.close()
  }
  send() {
    // TODO: add check for 4000 symbols
    const payload = {
      message: this.state.message
    }
  }
  render() {
    return (
      <div class="container">
          <div class="form-group">
            <label for="exampleInputEmail1">Let's broadcast!</label>
            <textarea rows="4" class="form-control" id="exampleInputEmail1"
              aria-describedby="emailHelp" placeholder="Enter message"
              onChange={this.handleChange} />
            <small id="emailHelp" class="form-text text-muted">
              Double check before sending! This action cant be undone.
              </small>
          </div>
          <button type="submit" class="btn btn-primary">Broadcast</button>
          <button class="btn btn-primary pull-right"
                  onCLick={this.canClose}>
                  Close</button>
      </div>
    )
  }
}

const UserInfo = ({ children, ...props }) => {
  let payload
  const msg = props.obj.last_message
  if (msg.text.length > 0) payload = msg.text.substring(0, 70)
  else if (msg.photo_id.length > 0) payload = "Photo"
  else payload = "File"
  const username = props.obj.username
  const name = props.obj.name
  const email = props.obj.email
  return (
    <userinfo>
      <div class={"chat_list" + (props.active ? " active_chat" : " ")} >
        <div class="chat_people">
          <div data-rh={email || "Not registered"} class="chat_img">
            <img src="/static/img/default/user.png" alt="sunil"> </img>
          </div>
          <div class="chat_ib" onClick={() => props.choose(props.obj.userid)}>
            <h5>{name}
                {username.length > 0 ? <div>{"@" + username}</div> : ""}
              <span class="chat_d">{formatDate(msg.send_date).day}</span>
            </h5>
            <p>{payload}</p>
          </div>
        </div>
      </div>
    </userinfo>
  )
}

class UserList extends Component {
  constructor(props) {
    super(props)
    this.setState({
      searchInput: "",
      userlist: props.userlist,
      broadcastModal: false,
    })
    this.handleChange = this.handleChange.bind(this)
    this.openModal = this.openModal.bind(this)
    this.closeModal = this.closeModal.bind(this)
    this.afterOpenModal = this.afterOpenModal.bind(this)
  }
  openModal() {
    this.setState({broadcastModal: true})
  }
  afterOpenModal() {}
  closeModal() {
    this.setState({broadcastModal: false})
  }
  handleChange(event) {
    this.setState({
      searchInput: event.target.value,
    })
  }
  render() {
    const search = (input, users) => {
      const fuse = new Fuse(users, fuseOptions)
      return fuse.search(input)
    }
    return (
      <div class="inbox_people">
        <div class="headind_srch">
          <div class="recent_heading" onClick={this.openModal}>
            <h4><a href="#">Broadcast</a></h4>
            <Modal
              isOpen={this.state.broadcastModal}
              onAfterOpen={this.afterOpenModal}
              onRequestClose={this.closeModal}
              shouldCloseOnOverlayClick={false}
              style={customStyles}
              contentLabel="Example modal">
              <BroadcastModal close={this.closeModal}/>
            </Modal>
          </div>
          <div class="srch_bar">
            <div class="stylish-input-group">
              <input type="text" class="search-bar" placeholder="Search..."
                value={this.state.searchInput} onInput={this.handleChange} />
            </div>
          </div>
        </div>
        {
          this.props.userlist.length > 0 ?
            search(this.state.searchInput, this.props.userlist).map(user =>
              <UserInfo obj={user} choose={this.props.choose}
                active={(this.props.current == user.userid)} />)
            : (<p>"No users"</p>)
        }
      </div>)
  }
}

class MessageForm extends Component {
  constructor(props) {
    super(props)
    this.state = {text: ''}

    this.handleChange = this.handleChange.bind(this)
    this.sendMessage = this.sendMessage.bind(this)
    this.handleFile = this.handleFile.bind(this)
  }
  handleChange(event) {
    this.setState({text: event.target.value})
  }
  handleFile(event) {

  }
  sendMessage() {
    console.log('sending message...')
    console.log(this.state)
    if (this.state.text.length == 0) return
    const data = {
      user_id: this.props.chatid,
      text: this.state.text,
    }
    const config = {
      method: 'post',
      url: `/api/user/${this.props.chatid}/messages`,
      data: data,
      headers: {
        'content-Type': 'application/json',
        'Accept': '/',
        'Cache-Control': 'no-cache',
      },
      credentials: 'same-origin',
    }
    axios(config)
    .then(response => {
      console.log(response)
      this.setState({
        text: ''
      })
    })
    .catch(err => console.log(err))
  }
  render() {
    return (
      <div class="type_msg">
        <div class="input_msg_write">
          <div class="input-group">
            <input type="text" class="write_msg" placeholder="Type a message"
                   value={this.state.text} onChange={this.handleChange} />
            <span class="input-group-btn">
              <button class="btn btn-default" onClick={() => this.sendMessage()}>
                Send
              </button>
            </span>
          </div>
        </div>
      </div>
    )
  }
}

class MessageData extends Component {
  constructor(props) {
    super(props)
    this.state = {
      modalIsOpen: false
    }
    this.openModal = this.openModal.bind(this)
    this.afterOpenModal = this.afterOpenModal.bind(this)
    this.closeModal = this.closeModal.bind(this)
  }
  openModal() {
    this.setState({modalIsOpen: true})
  }
  afterOpenModal() {

  }
  closeModal() {
    this.setState({modalIsOpen: false})
  }
  render() {
    let payload
    const props = this.props
    if (props.obj.text.length > 0) {
      payload = (<p>{props.obj.text}</p>)
    } else if (props.obj.photo_id.length > 0) {
      payload = (
        <div>
          <div onClick={this.openModal}>
            <img src={`/file/${props.obj.photo_id}`} />
          </div>
          <Modal
            isOpen={this.state.modalIsOpen}
            onAfterOpen={this.afterOpenModal}
            onRequestClose={this.closeModal}
            style={customStyles}
            contentLabel="Example Modal" >
            <img src={`/file/${props.obj.photo_id}`} />
          </Modal>
        </div>
    )
    } else {
      payload = (<a href={`/file/${props.obj.file_id}`}>File</a>)
    }
    return (
      <div>
        {payload}
        <span class="time_date">{formatDate(props.obj.send_date).time}| {formatDate(props.obj.send_date).day}</span>
      </div>
    )
  }
}

const IncomingMessage = ({ children, ...props }) => (
  <div class="incoming_msg">
    <div  class="incoming_msg_img">
      <img src="/static/img/default/user.png" alt="sunil">
      </img>
    </div>
    <div class="received_msg">
      <div class="received_withd_msg">
        <MessageData obj={props.obj} />
      </div>
    </div>
  </div>
)

const SentMessage = ({ children, ...props }) => (
  <div class="outgoing_msg">
    <div class="sent_msg">
      <MessageData obj={props.obj} />
    </div>
  </div>
)

const MessageHistory = ({ children, ...props }) => (
  <div class="mesgs">
    <div class="msg_history" >
      {
        props.messages ?
        props.messages.map(msg => {
          if (msg.from_bot) return (<SentMessage obj={msg} />)
          return (<IncomingMessage obj={msg} />)
        })
        :
        <p>No messages yet :(</p>
      }
    </div>
    <MessageForm chatid={props.chatid}/>
  </div>
)

class App extends Component {
  constructor() {
    super()
    this.state.users = {}
    this.state.messages = {}
    this.state.currentChat = false
    this.state.websocket = {}
    this.chatClicked = this.chatClicked.bind(this)
    this.processUpdate = this.processUpdate.bind(this)
  }
  chooseChat(id) {
    // TODO: cache messages to app state
    axios.get(`/api/user/${id}/messages`)
    .then(response => {
      let updated = Object.assign({[id]: response.data}, this.state.messages)
      this.setState({
        messages: updated,
        currentChat: id
      })
    })
    .catch(err => {
      console.log(err)
    })
  }
  processUpdate(event) {
    console.log(event)
    let msg = JSON.parse(event.data)
    if (msg.chat_id == this.state.currentChat) {
      let history = this.state.messages[this.state.currentChat]
      history.unshift(msg)
      let updated = Object.assign({[msg.chat_id]: history}, this.state.messages)
      console.log(updated)
      this.setState({
        messages: updated,
      })
    }
  }
  componentDidMount() {
    // ws request would use document cookies, we're safe
    // receive ws connection, pass it to children
    const href = document.location.host
    var ws
    if (document.location.protocol === "https:") {
      ws = new WebSocket(`wss://${href}/ws`)
    } else {
      ws = new WebSocket(`ws://${href}/ws`)
    }
    ws.onopen = function(event) {
      console.log('Socket opened!')
    }
    ws.onmessage = (event) => {
      // if received message chat is equal to current chat, append message to history
      // if not, and message from user (NOT admin) set message's chat as unread
      // Update last message in user list too.
      this.processUpdate(event)
    }
    ws.onerror = function(e) {
      console.log(e)
    }
    ws.onclose = function(e) {
      console.log(e)
    }
    // because I dont care
    console.log(ws)
    this.setState({
      websocket: ws
    })
    // get user list (backend sort it for us)
    axios.get('/api/user/', axiosConfig)
    .then(response => {
      console.log(response)
      this.setState({
        users: response.data
      })
      const latestUser = response.data[0].userid
      this.chooseChat(latestUser)
    }).catch(err => {
      console.log(err)
    })
    // load message history for first chat in list
  }
  currentChatHistory() {
    if (this.state.currentChat === false) return new Array()
    return this.state.messages[this.state.currentChat]
  }
  chatClicked(id) {
    console.log('chat chosen: ', id)
    // TODO: set loader
    this.chooseChat(id)
  }
	render({}, { text }) {
		return (
	 		<app>
        <ReactHint autoPosition events delay={{show: 100, hide: 1000}} />
        <div class="messaging" >
          <div class="inbox_msg">
          <UserList userlist={this.state.users}
                    choose={this.chatClicked}
                    current={this.state.currentChat}/>
          <MessageHistory messages={this.currentChatHistory()}
                    chatid={this.state.currentChat}/>
          </div>
        </div>
			</app>
		);
	}
}

// Start 'er up:
render(<App />, document.getElementById('chat'));
