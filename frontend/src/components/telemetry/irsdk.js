function worker_func() {
    console.log("worker")
}

export default worker_func;

if (window!=self) {
    worker_func();
}