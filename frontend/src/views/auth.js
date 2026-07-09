import { validateEmail, validatePassword } from '../utils/validators.js';
import { login, register, getProfile } from '../api/auth.js';

// Элементы DOM
const authForm = document.getElementById('auth-form');
const authTitle = document.getElementById('auth-title');
const submitBtn = document.getElementById('submit-btn');
const usernameGroup = document.getElementById('username-group');
const usernameInput = document.getElementById('username');
const emailInput = document.getElementById('email');
const passwordInput = document.getElementById('password');
const errorContainer = document.getElementById('error-message');
const toggleLink = document.getElementById('toggle-mode-link');

// Состояние экрана (true = Вход, false = Регистрация)
let isLoginMode = true;

// 1. Логика переключения режимов (Вход / Регистрация)
toggleLink.addEventListener('click', (event) => {
    event.preventDefault(); // Отменяем переход по ссылке
    isLoginMode = !isLoginMode; // Меняем флаг

    // Очищаем ошибки прошлых попыток
    hideError();

    if (isLoginMode) {
        authTitle.textContent = 'Войти в аккаунт';
        submitBtn.textContent = 'Войти';
        toggleLink.textContent = 'Нет аккаунта? Зарегистрироваться';
        usernameGroup.classList.add('hidden');
        usernameInput.removeAttribute('required');
    } else {
        authTitle.textContent = 'Регистрация';
        submitBtn.textContent = 'Создать аккаунт';
        toggleLink.textContent = 'Уже есть аккаунт? Войти';
        usernameGroup.classList.remove('hidden');
        usernameInput.setAttribute('required', 'true');
    }
});

// 2. Обработка отправки формы
authForm.addEventListener('submit', async (event) => {
    event.preventDefault(); // КРИТИЧНО: предотвращаем ту самую перезагрузку с "?"
    hideError();

    const email = emailInput.value.trim();
    const password = passwordInput.value;
    const username = usernameInput.value.trim();

    // Фронтенд валидация пароля
    if (!validatePassword(password)) {
        showError('Пароль должен быть не менее 6 символов');
        return;
    }

    // Если мы в режиме регистрации, проверим еще и имя
    if (!isLoginMode && username.length < 2) {
        showError('Имя пользователя должно быть не менее 2 символов');
        return;
    }

    try {
        let response;

        // Вместо сырого fetch вызываем твои готовые функции из api/auth.js
        if (isLoginMode) {
            response = await login(email, password);
        } else {
            response = await register(email, password, username);
        }

        // Твой apiRequest возвращает оригинальный объект response,
        // поэтому вся дальнейшая логика проверки статусов остается прежней!
        if (response.ok) {
            // Успех — летим на дашборд
            window.location.href = '/views/dashboard.html';
        } else {
            // Ошибка от бэкенда (400, 401, 409)
            const errorData = await response.json().catch(() => ({}));
            const message = errorData.message || 'Что-то пошло не так. Попробуйте снова.';
            showError(message);
        }
    } catch (err) {
        console.error('Сетевая ошибка:', err);
        showError('Не удалось связаться с сервером. Проверьте сеть.');
    }
});

// Вспомогательные функции для отображения ошибок
function showError(text) {
    errorContainer.textContent = text;
    errorContainer.classList.remove('hidden');
}

function hideError() {
    errorContainer.textContent = '';
    errorContainer.classList.add('hidden');
}
async function checkExistingAuth() {
    // Важно: временно отключаем авто-редирект из client.js для этого вызова,
    // но так как мы проверяем pathname внутри client.js (auth.html), он и так не сработает.
    try {
        const response = await getProfile();
        if (response.ok) {
            // Если сессия жива — пользователю нечего делать на странице логина
            window.location.href = '/views/dashboard.html';
        }
    } catch (error) {
        // Если ошибка — всё отлично, значит пользователь не авторизован и должен видеть форму
        console.log('Активной сессии не найдено, показываем форму входа.');
    }
}

// Запускаем проверку
checkExistingAuth();