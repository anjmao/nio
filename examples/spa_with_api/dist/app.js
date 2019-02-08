function fetchData() {
    fetch('/api/data')
        .then((res) => res.json())
        .then(render)
}

function render(data) {
    const root = document.querySelector('#app');
    for (let item of data) {
        const div = document.createElement('div');
        div.innerText = item.name;
        root.appendChild(div)
    }
}