async function sendExpression() {
    const expression = document.getElementById("expression").value;
    const resultField = document.getElementById("result");
    const errorField = document.getElementById("error");

    if (!expression) {
        alert("Введите выражение!");
        return;
    }

    // Очистка полей перед началом вычисления
    resultField.innerText = "Вычисляется...";
    errorField.innerText = "";

    try {
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
            // Вывод ошибки в отдельное поле
            errorField.innerText = data.error || "Неизвестная ошибка";
            resultField.innerText = "—"; // Сбрасываем поле результата
        }
    } catch (error) {
        errorField.innerText = "Ошибка запроса к серверу";
        resultField.innerText = "—"; // Сбрасываем поле результата
    }
}

// Функция для получения результата с повторными запросами
async function checkResult(exprID) {
    const resultField = document.getElementById("result");
    const errorField = document.getElementById("error");
    let attempts = 30;  // Количество попыток (по 3 сек = 90 сек макс)

    while (attempts > 0) {
        await new Promise(r => setTimeout(r, 3000)); // Ждём 3 секунды

        try {
            const response = await fetch(`http://localhost:8080/api/v1/result/${exprID}`);
            if (response.ok) {
                const data = await response.json();
                resultField.innerText = data.result; // Выводим только результат
                return;
            }
        } catch (error) {
            errorField.innerText = "Ошибка получения результата";
        }

        attempts--;
    }

    errorField.innerText = "Ошибка: результат не получен";
    resultField.innerText = "—"; // Сбрасываем поле результата
}
