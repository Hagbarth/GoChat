var globalState = {}
$(function () {
  // Setup websocket
  var ws = new WebSocket('ws://localhost:8000/socketserver')
  ws.onmessage = function (e) {
    var messages = JSON.parse(e.data)
    updateState(globalState, {
      messages: messages
    })
  }

  function updateState (state, added) {
    globalState = $.extend(state, added)
    render(globalState)
  }

  function bindAction (fn, type) {
    return function (state) {
      return fn.apply(null, arguments)
      .then(function (r) {
        updateState(state, {
          [type]: r
        })
      })
    }
  }

  function callApi (path, body, method) {
    return fetch(path, {
      method: method || 'GET',
      headers: {
        'content-type': 'application/json'
      },
      body: JSON.stringify(body)
    })
    .then(function (r) {
      return r.json()
    })
  }

  var login = bindAction(function(state, nickname) {
    return callApi('/login', {
      nickname: nickname
    }, 'POST')
  }, 'user')

  var sendMessage = function(state, user, message) {
    ws.send(JSON.stringify({
      Uid: user.ID,
      message: message
    }))
  }

  var getMessages = bindAction(function(state) {
    return callApi('/messages')
  }, 'messages')

  function render(state) {
    console.log(Object.assign({}, state))
    var $loginForm = $('#login-form')
    var $loginFormSubmit = $('#login-form-submit')
    var $chatboxContainer = $('#chatbox-container')
    var $chatbutton = $('#chatbutton')
    var $chatbox = $('#chatbox')
    var $messageboard = $('#messageboard')

    /**
    * Handle UI
    **/
    if (state.user) {
      $loginForm.hide()
      $messageboard.show()
      $chatboxContainer.show()
    } else {
      $loginForm.show()
      $chatboxContainer.hide()
      $messageboard.hide()
    }

    if (state.messages) {
      var messagesHTML = ''
      for (var i = 0; i < state.messages.length; i++) {
        var message = state.messages[i]
        messagesHTML += '<div class="chat-message"><p class="sender">' + message.User.Nickname + '</p>' + message.Message + '</div>'
      }
      $messageboard.html(messagesHTML)
    }

    /**
    * Handle events
    **/
    // Login
    $loginForm.off()
    $loginForm.on('submit', function(e) {
      e.preventDefault()
      var nickname = $(this).find('input[name="nickname"]').val()
      login(state, nickname)
      .then(function () {
        getMessages(state)
      })
    })
    $loginFormSubmit.off()
    $loginFormSubmit.on('click', function() {
      $loginForm.submit()
    })

    // Chat
    $chatbutton.off()
    $chatbutton.on('click', function () {
      var message = $chatbox.val()
      $chatbox.val('')
      sendMessage(state, state.user, message)
    })
  }
  render({})
})
