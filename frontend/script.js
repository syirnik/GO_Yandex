async function sendExpression() {
    const expression = document.getElementById("expression").value;
    if (!expression) {
        alert("Введите выражение!");
        return;
    }

    document.getElementById("result").innerText = "Вычисляется...";

    const response = await fetch("http://localhost:8080/api/v1/calculate", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ expression: expression })
    });

    const data = await response.json();

    if (response.ok) {
        await checkResult(data.id);
    } else {
        document.getElementById("result").innerText = `Ошибка: ${data.error || "Неизвестная ошибка"}`;
    }
}


// Функция для получения результата с повторными запросами
async function checkResult(exprID) {
    let attempts = 30;  // Количество попыток (по 3 сек = 60 сек макс)
    while (attempts > 0) {
        await new Promise(r => setTimeout(r, 3000)); // Ждём 3 секунды

        const response = await fetch(`http://localhost:8080/api/v1/result/${exprID}`);
        if (response.ok) {
            const data = await response.json();
            document.getElementById("result").textContent = `${data.result}`;
            return; // Выходим из функции, чтобы не обновлять результат ещё раз
        }

        attempts--;
    }
    document.getElementById("result").innerText = "Ошибка: результат не получен";
}
