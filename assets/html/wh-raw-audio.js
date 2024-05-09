class WHRawAudioElement extends HTMLElement {
  static observedAttributes = ["src", "id"];
  constructor() {
    super();
    this.btnStatus = "Pause";
    const shadowRoot = this.attachShadow({ mode: "closed" });
    this.audio = document.createElement("audio");
    this.btn = document.createElement("button");
    this.btn.textContent = "播放原文";
    this.btn.onclick = () => {
      if (this.btnStatus == "Play") {
        this.btnStatus = "Pause";
        this.btn.textContent = "播放原文";
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
    this.btn.textContent = "播放原文";
  }
}

customElements.define("wh-raw-audio", WHRawAudioElement);
