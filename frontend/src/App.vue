<template>
	<div class="chatbox">
		<h2>Anonymous Chat App</h2>
		<div v-for="message in messages" :key="message.msg">
			<div v-if="message.user == this.user">
				<!-- Add CSS -->
				<p>{{ message.User }}: {{ message.Text }}</p>
			</div>
			<div v-else>
				<!-- Add CSS -->
				<p>{{ message.User }}: {{ message.Text }}</p>
			</div>
		</div>
		<div>
			<input v-model="message" type="text" />
			<button v-on:click="sendMessage">Submit</button>
		</div>
	</div>
</template>

<script>
export default {
	name: 'App',
	data() {
		return {
			message: '',
			messages: [],
			user: '',
			socket: null,
		};
	},
	mounted() {
		if (this.socket == null) {
			this.user = Math.random().toString(36).substring(2, 7);
			this.socket = new WebSocket('ws://localhost:4500/socket/' + this.user);
			this.socket.onmessage = msg => this.acceptMessage(msg);
		}
	},
	methods: {
		sendMessage() {
			let msg = {
				Text: this.message,
				User: this.user,
			};
			this.socket.send(JSON.stringify(msg));
		},
		acceptMessage(msg) {
			console.log(msg);
			this.messages.push(JSON.parse(msg.data));
		},
	},
};
</script>

<style>
#app {
	font-family: Avenir, Helvetica, Arial, sans-serif;
	-webkit-font-smoothing: antialiased;
	-moz-osx-font-smoothing: grayscale;
	text-align: center;
	color: #2c3e50;
	margin-top: 60px;
}
body {
	padding: 0;
	margin: 0;
	display: flex;
	justify-content: center;
}
.chatbox {
	width: 350px;
	border: 12px solid rgb(50, 158, 50);
	border-radius: 15px;
	padding-right: 10px;
	padding-left: 10px;
	padding-top: 5px;
	padding-bottom: 5px;
}
#message-container {
	width: 80%;
	max-width: 1200px;
}

#message-container div {
	background-color: #ccc;
	padding: 5px;
}

#message-container div:nth-child(2n) {
	background-color: #fff;
}

#send-container {
	position: fixed;
	padding-bottom: 30px;
	bottom: 0;
	background-color: white;
	max-width: 1200px;
	width: 80%;
	display: flex;
}

#message-input {
	flex-grow: 1;
}
</style>
