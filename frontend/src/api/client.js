// Базовый URL нашего API. Когда поднимем Nginx, запросы пойдут на него.
const BASE_URL = '/api';

export async function apiRequest(endpoint, options = {}) {
    const url = `${BASE_URL}${endpoint}`;
    
    // Задаем дефолтные заголовки, если они не переданы
    options.headers = {
        'Content-Type': 'application/json',
        ...options.headers,
    };

    try {
        const response = await fetch(url, options);
        
        // Перехватываем глобальные системные штуки (например, если токен протух)
        if (response.status === 401) {
            console.warn('Пользователь не авторизован (401)');
    
            // Проверяем, что мы сейчас НЕ на странице auth.html
            if (!window.location.pathname.includes('auth.html')) {
                // Только в этом случае делаем жесткий редирект на логин
                window.location.href = '/views/auth.html';
            }
        }

        return response;
    } catch (error) {
        console.error(`Сетевая ошибка при запросе к ${endpoint}:`, error);
        throw error; // Прокидываем ошибку дальше в хендлер страницы
    }
}