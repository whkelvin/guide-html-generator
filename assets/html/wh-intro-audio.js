class WHIntroAudioElement extends HTMLElement {
  static observedAttributes = ["src", "id"];
  constructor() {
    super();
    this.btnStatus = "Pause";
    const shadowRoot = this.attachShadow({ mode: "closed" });
    this.audio = document.createElement("audio");
    this.audio.onended = () => {
      this.audioEnded();
    };
    this.audio.onerror = () => {
      this.audio.setAttribute("src", "./assets/mp3/default.mp3");
      this.audio.load();
    };
    this.btn = document.createElement("button");
    this.btn.textContent = "播放簡介";
    this.btn.onclick = () => {
      if (this.btnStatus == "Play") {
        this.btnStatus = "Pause";
        this.btn.textContent = "播放簡介";
        this.audio.pause();
      } else {
        this.btnStatus = "Play";
        this.btn.textContent = "暫停";
        this.audio.play();
      }
    };
    console.log(this.btn.textContent);
    shadowRoot.append(this.audio);
    shadowRoot.append(this.btn);
  }

  attributeChangedCallback(name, oldValue, newValue) {
    console.log(
      `Attribute ${name} has changed from ${oldValue} to ${newValue}.`
    );
    this.audio.setAttribute("src", newValue);
    this.audio.load();
    this.btnStatus = "Pause";
    this.btn.textContent = "播放簡介";
  }

  audioEnded() {
    this.audio.load();
    this.btnStatus = "Pause";
    this.btn.textContent = "播放簡介";
  }
}

customElements.define("wh-intro-audio", WHIntroAudioElement);
