class RT {
  ws = null;
  connect(url) {
    try { this.ws?.close(); } catch {}
    this.ws = new WebSocket(url);
  }
  disconnect() {
    try { this.ws?.close(); } catch {}
    this.ws = null;
  }
}
const realtime = new RT();
export default realtime;
