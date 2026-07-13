import { validateEmail, validatePassword } from '../utils/validators.js';
import { login, register, getProfile, verify2FA } from '../api/auth.js';

// Элементы DOM
const authForm = document.getElementById('auth-form');
const authTitle = document.getElementById('auth-title');
const submitBtn = document.getElementById('submit-btn');
const usernameGroup = document.getElementById('username-group');
const usernameInput = document.getElementById('username');
const emailInput = document.getElementById('email');
const emailGroup = emailInput.closest('.form-group');
const passwordInput = document.getElementById('password');
const passwordGroup = passwordInput.closest('.form-group');
const errorContainer = document.getElementById('error-message');
const toggleLink = document.getElementById('toggle-mode-link');
const verify2faGroup = document.getElementById('2fa-group');
const verify2faInput = document.getElementById('2fa-code');

// Состояние экрана: 'login' | 'register' | '2fa'
let mode = 'login';
// Email, для которого запрошен код — нужен на шаге verify2fa,
// т.к. сам email-инпут на этом шаге скрыт/задизейблен
let pendingEmail = '';

// 1. Отрисовка состояния формы под текущий режим
function renderMode() {
    hideMessage();

    if (mode === 'login') {
        authTitle.textContent = 'Войти в аккаунт';
        submitBtn.textContent = 'Продолжить';
        toggleLink.textContent = 'Нет аккаунта? Зарегистрироваться';

        usernameGroup.classList.add('hidden');
        usernameInput.removeAttribute('required');

        verify2faGroup.classList.add('hidden');
        verify2faInput.removeAttribute('required');

        emailGroup.classList.remove('hidden');
        passwordGroup.classList.remove('hidden');
        emailInput.removeAttribute('disabled');
        passwordInput.removeAttribute('disabled');
    } else if (mode === 'register') {
        authTitle.textContent = 'Регистрация';
        submitBtn.textContent = 'Создать аккаунт';
        toggleLink.textContent = 'Уже есть аккаунт? Войти';

        usernameGroup.classList.remove('hidden');
        usernameInput.setAttribute('required', 'true');

        verify2faGroup.classList.add('hidden');
        verify2faInput.removeAttribute('required');

        emailGroup.classList.remove('hidden');
        passwordGroup.classList.remove('hidden');
        emailInput.removeAttribute('disabled');
        passwordInput.removeAttribute('disabled');
    } else if (mode === '2fa') {
        authTitle.textContent = 'Подтверждение входа';
        submitBtn.textContent = 'Подтвердить';
        toggleLink.textContent = 'Отмена, ввести данные заново';

        usernameGroup.classList.add('hidden');
        usernameInput.removeAttribute('required');

        // Email и пароль на этом шаге больше не нужны — прячем и блокируем,
        // чтобы форма их не отправляла и не путала пользователя
        emailGroup.classList.add('hidden');
        passwordGroup.classList.add('hidden');
        emailInput.setAttribute('disabled', 'true');
        passwordInput.setAttribute('disabled', 'true');

        verify2faGroup.classList.remove('hidden');
        verify2faInput.setAttribute('required', 'true');
        verify2faInput.value = '';
        verify2faInput.focus();
    }
}

// 2. Переключение режимов по ссылке
toggleLink.addEventListener('click', (event) => {
    event.preventDefault();

    if (mode === '2fa') {
        // Отмена верификации — возвращаемся к чистой форме логина
        mode = 'login';
        pendingEmail = '';
    } else {
        mode = mode === 'login' ? 'register' : 'login';
    }

    renderMode();
});

// 3. Обработка отправки формы
authForm.addEventListener('submit', async (event) => {
    event.preventDefault(); // предотвращаем перезагрузку страницы
    hideMessage();

    if (mode === '2fa') {
        await handleVerify2FA();
        return;
    }

    const email = emailInput.value.trim();
    const password = passwordInput.value;
    const username = usernameInput.value.trim();

    if (!validateEmail(email)) {
        showError('Введите корректный email');
        return;
    }
    if (!validatePassword(password)) {
        showError('Пароль должен быть не менее 6 символов');
        return;
    }
    if (mode === 'register' && username.length < 2) {
        showError('Имя пользователя должно быть не менее 2 символов');
        return;
    }

    if (mode === 'login') {
        await handleLogin(email, password);
    } else {
        await handleRegister(email, password, username);
    }
});

// --- Шаг 1: логин + пароль ---
async function handleLogin(email, password) {
    try {
        const response = await login(email, password);

        if (response.ok) {
            const data = await response.json().catch(() => ({}));
            if (data.status) {
                // Бэкенд подтвердил учётные данные и отправил код — переходим к вводу кода
                pendingEmail = email;
                mode = '2fa';
                renderMode();
            } else {
                // На случай, если бэкенд когда-нибудь начнёт пускать без 2FA
                window.location.href = '/views/dashboard.html';
            }
        } else if (response.status === 401) {
            showError('Неверный email или пароль');
        } else {
            showError('Что-то пошло не так. Попробуйте снова.');
        }
    } catch (err) {
        console.error('Сетевая ошибка:', err);
        showError('Не удалось связаться с сервером. Проверьте сеть.');
    }
}

// --- Шаг 2: ввод кода из письма ---
async function handleVerify2FA() {
    const code = verify2faInput.value.trim();

    if (!/^\d{6}$/.test(code)) {
        showError('Код должен состоять из 6 цифр');
        return;
    }

    try {
        const response = await verify2FA(pendingEmail, code);

        if (response.ok) {
            // Сессионная кука уже выставлена бэкендом — летим на дашборд
            window.location.href = '/views/dashboard.html';
        } else if (response.status === 400) {
            const errorData = await response.json().catch(() => ({}));
            showError(errorData.error || 'Неверный или истекший код подтверждения');
        } else {
            showError('Что-то пошло не так. Попробуйте снова.');
        }
    } catch (err) {
        console.error('Сетевая ошибка:', err);
        showError('Не удалось связаться с сервером. Проверьте сеть.');
    }
}

// --- Регистрация ---
async function handleRegister(email, password, username) {
    try {
        const response = await register(email, password, username);

        if (response.status === 201) {
            // После регистрации отдельного шага логина всё равно требует 2FA —
            // просто переключаем пользователя на форму входа
            mode = 'login';
            renderMode();
            showSuccess('Регистрация успешна. Теперь войдите в аккаунт.');
        } else if (response.status === 409) {
            showError('Пользователь с таким email уже существует');
        } else {
            showError('Что-то пошло не так. Попробуйте снова.');
        }
    } catch (err) {
        console.error('Сетевая ошибка:', err);
        showError('Не удалось связаться с сервером. Проверьте сеть.');
    }
}

// Вспомогательные функции для отображения сообщений
function showError(text) {
    errorContainer.textContent = text;
    errorContainer.classList.remove('hidden', 'success-text');
}

function showSuccess(text) {
    errorContainer.textContent = text;
    errorContainer.classList.remove('hidden');
    errorContainer.classList.add('success-text');
}

function hideMessage() {
    errorContainer.textContent = '';
    errorContainer.classList.add('hidden');
    errorContainer.classList.remove('success-text');
}

async function checkExistingAuth() {
    try {
        const response = await getProfile();
        if (response.ok) {
            // Если сессия жива — пользователю нечего делать на странице логина
            window.location.href = '/views/dashboard.html';
        }
    } catch (error) {
        console.log('Активной сессии не найдено, показываем форму входа.');
    }
}

// Запускаем проверку и первичную отрисовку
renderMode();
checkExistingAuth();