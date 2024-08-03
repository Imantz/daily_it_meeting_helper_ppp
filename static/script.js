const ws = new WebSocket('ws://localhost:8080/ws');

const progressInput = document.getElementById('progress');
const plansInput = document.getElementById('plans');
const problemsInput = document.getElementById('problems');
const generateButton = document.getElementById('generate');
const responseDiv = document.getElementById('response');

ws.onopen = () => {
    console.log('WebSocket connection established');
};

ws.onmessage = (event) => {
    console.log('Received message:', event.data);
};

const sendData = () => {
    const data = {
        progress: progressInput.value,
        plans: plansInput.value,
        problems: problemsInput.value
    };
    ws.send(JSON.stringify(data));
};

progressInput.addEventListener('input', sendData);
plansInput.addEventListener('input', sendData);
problemsInput.addEventListener('input', sendData);

generateButton.addEventListener('click', async () => {
    const data = {
        progress: progressInput.value,
        plans: plansInput.value,
        problems: problemsInput.value
    };

    const response = await fetch('http://localhost:8080/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });

    const result = await response.json();
    responseDiv.textContent = result.formattedText;
});

window.addEventListener('load', async () => {
    const response = await fetch('http://localhost:8080/current-entry');
    const data = await response.json();
    console.log(data);
    progressInput.value = data.progress || '';
    plansInput.value = data.plans || '';
    problemsInput.value = data.problems || '';
});
