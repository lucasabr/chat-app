<template>
	<div id="main" class="chatbox">
		<h2>Anonymous Chat App</h2>
		<div v-for="message in messages" :key="message.msg">
			<div class="message current" v-if="message.User == this.user">
				<!-- Add CSS -->
				<p class="user">{{ message.User }} (You)</p>
				<p v-if="message.Type == 'message'">{{ message.Text }}</p>
				<audio :src="message.Chunks" v-else controls></audio>
			</div>
			<div class="message other" v-else>
				<!-- Add CSS -->
				<p class="user">{{ message.User }}</p>
				<p v-if="message.Type == 'message'">{{ message.Text }}</p>
				<audio :src="message.Chunks" v-else controls></audio>
			</div>
		</div>
		<div>
			<input v-model="message" type="text" />
			<button v-on:click="sendMessage">Submit</button>
		</div>
		<div>
			<button v-on:click="toggleRec">Start Recording</button>
			<button v-on:click="sendAudioMessage">Send</button>
		</div>
		<audio id="audio-preview" controls></audio>
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
			mediaRec: null,
			chunks: [],
		};
	},
	mounted() {
		if (this.socket == null) {
			this.user = Math.random().toString(36).substring(2, 7);
			this.socket = new WebSocket('ws://localhost:4500/socket/' + this.user);
			this.socket.onmessage = msg => this.acceptMessage(msg);
		}
		if (navigator.mediaDevices && navigator.mediaDevices.getUserMedia) {
			navigator.mediaDevices
				.getUserMedia({
					audio: true,
				})
				.then(stream => {
					if (this.mediaRec == null) {
						this.mediaRec = new MediaRecorder(stream);
						this.mediaRec.ondataavailable = e => this.chunks.push(e.data);
						this.mediaRec.onstop = () => {
							const blob = new Blob(this.chunks, { type: 'audio/ogg; codecs=opus' });
							const audioURL = window.URL.createObjectURL(blob);
							document.getElementById('audio-preview').src = audioURL;
						};
					}
				})
				.catch(err => {
					console.log(err);
				});
		}
	},
	methods: {
		sendMessage() {
			let msg = {
				Type: 'message',
				Text: this.message,
				User: this.user,
				Chunks: '',
			};
			this.socket.send(JSON.stringify(msg));
		},
		sendAudioMessage() {
			if (this.chunks.length == 0) return;
			let msg = {
				Type: 'audio',
				Text: '',
				User: this.user,
				Chunks: document.getElementById('audio-preview').src,
			};
			this.socket.send(JSON.stringify(msg));
			this.chunks = [];
			document.getElementById('audio-preview').src = '';
		},
		acceptMessage(msg) {
			console.log(msg);
			this.messages.push(JSON.parse(msg.data));
		},
		toggleRec(event) {
			if (this.mediaRec.state == 'inactive') {
				this.mediaRec.start();
				event.target.style.background = 'red';
				event.target.innerText = 'Stop Recording';
			} else {
				this.mediaRec.stop();
				event.target.style.background = 'green';
				event.target.innerText = 'Start Recording';
			}
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
	border: 12px solid rgb(50, 158, 50);
	border-radius: 15px;
	padding-right: 10px;
	padding-left: 10px;
	padding-top: 5px;
	padding-bottom: 5px;
}
.message {
	border-radius: 15px;
	border: 2px solid white;
	background-color: rgb(98, 93, 116);
}
.message p {
	color: white;
}

.user {
	font-size: 10px;
}
</style>
