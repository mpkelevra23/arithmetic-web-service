// Выделяем Vue объекты для работы
const { createApp, ref } = Vue;

createApp({
    setup() {
        const expression = ref('');
        const result = ref(null);
        const error = ref('');
        const isCalculating = ref(false);

        const calculate = async () => {
            // Проверка на пустую строку
            if (!expression.value.trim()) {
                error.value = 'Please enter an expression';
                return;
            }

            error.value = '';
            isCalculating.value = true;

            try {
                const response = await fetch('/api/v1/calculate', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ expression: expression.value }),
                });

                const data = await response.json();

                if (!response.ok) {
                    error.value = data.error || 'An error occurred';
                    result.value = null;
                } else {
                    result.value = data.result;
                }
            } catch (err) {
                error.value = 'Failed to communicate with the server';
                console.error('Error:', err);
            } finally {
                isCalculating.value = false;
            }
        };

        return {
            expression,
            result,
            error,
            isCalculating,
            calculate
        };
    }
}).mount('#app');