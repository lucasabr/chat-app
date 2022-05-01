<template>
	<div id="main" class="chatbox">
		<h2>Anonymous Chat App</h2>
		<div v-for="message in messages" :key="message.msg">
			<div v-if="message.User == this.user">
				<!-- Add CSS -->
				<p v-if="message.Type == 'message'">{{ message.User }}: {{ message.Text }}</p>
				<audio :src="message.Chunks" v-else controls></audio>
			</div>
			<div v-else>
				<!-- Add CSS -->
				<p v-if="message.Type == 'message'">{{ message.User }}: {{ message.Text }}</p>
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
							const audio = document.createElement('audio');
							const blob = new Blob(this.chunks, { type: 'audio/ogg; codecs=opus' });
							const audioURL = window.URL.createObjectURL(blob);
							audio.src = audioURL;
							audio.setAttribute('controls', '');
							audio.setAttribute('id', 'audioPreview');
							document.getElementById('main').appendChild(audio);
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
			const blob = new Blob(this.chunks, { type: 'audio/ogg; codecs=opus' });
			const audioURL = window.URL.createObjectURL(blob);
			let msg = {
				Type: 'audio',
				Text: '',
				User: this.user,
				Chunks: audioURL,
			};
			this.socket.send(JSON.stringify(msg));
			this.chunks = [];
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
