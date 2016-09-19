$(function () {
	function updateState (state, added) {
		render($.extend(true, state, added))
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

	var sendMessage = bindAction(function(state, user, message) {
		return callApi('/messages', {
			Uid: user.ID,
			message: message
		}, 'POST')
	}, 'messages')

	var getMessages = bindAction(function(state) {
		return callApi('/messages')
	}, 'messages')
	
	function render(state) {
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
			$loginForm.addClass('hidden')
			$chatboxContainer.removeClass('hidden')
		} else {
			$loginForm.removeClass('hidden')
			$chatboxContainer.addClass('hidden')
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