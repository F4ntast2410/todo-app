import { changeInfo, changeEmailVerify, changePassword, checkLinkTg, linkTg, setEmail, setEmailPassword, setEmailVerify2FA } from "../api/personal_account.js";
import { setupDashboardUI } from "./dashboard.js";
async function initTelegramSection(showMsg) {
    const container = document.getElementById('tg-widget-container');
    if (!container) return;

    try {
        // 1. Запрашиваем у бэкенда текущий статус привязки
        const response = await checkLinkTg();
        if (!response.ok) throw new Error('Не удалось получить статус привязки');
        
        const data = await response.json();
        if (data.linked) {
            // 2. Если уже привязан — рендерим статус
            renderTelegramLinked(container, data.username);
        } else {
            // 3. Если не привязан — рендерим кнопку виджета
            renderTelegramWidget(container, showMsg);
        }
    } catch (err) {
        container.innerHTML = `<p class="error-text">Не удалось загрузить данные интеграции</p>`;
    }
}

// Рендер состояния "Аккаунт привязан"
function renderTelegramLinked(container, username) {
    container.innerHTML = `
        <div class="tg-linked-badge">
            <span class="tg-icon">TG</span>
            <div class="tg-info">
                <p class="tg-title">Telegram привязан</p>
                <p class="tg-username">@${username}</p>
            </div>
        </div>
    `;
}

// Рендер самого виджета Telegram
function renderTelegramWidget(container, showMsg) {
    container.innerHTML = `<div id="tg-login-button"></div>`;

    // Создаем глобальный коллбэк, который вызовет скрипт Telegram при успешном входе
    window.onTelegramAuth = async function(user) {
        try {
            // user содержит: id, first_name, last_name, username, photo_url, auth_date, hash
            // Отправляем все эти поля на твой POST эндпоинт
            const response = await linkTg(user);

            if (response.ok) {
                showMsg('Telegram успешно привязан!', true);
                renderTelegramLinked(container, user.username);
            } else if (response.status === 409) {
                // Твой обработчик конфликта (StatusConflict) на бэкенде
                const errData = await response.json();
                showMsg(errData.error || 'Этот Telegram-аккаунт уже занят другим пользователем');
            } else {
                showMsg('Ошибка авторизации Telegram. Попробуйте еще раз.');
            }
        } catch (err) {
            showMsg('Ошибка сети при привязке Telegram');
        }
    };

    // Динамически загружаем скрипт виджета Telegram
    const script = document.createElement('script');
    script.src = "https://telegram.org/js/telegram-widget.js?22";
    script.async = true;
    
    // Настройки виджета
    script.setAttribute('data-telegram-login', 'todo_test42_bot'); // Замени на имя своего бота
    script.setAttribute('data-size', 'large');
    script.setAttribute('data-radius', '10');
    script.setAttribute('data-onauth', 'onTelegramAuth(user)'); // Указываем наш глобальный коллбэк
    script.setAttribute('data-request-access', 'write');

    // Вставляем скрипт, чтобы Telegram отрендерил кнопку внутри контейнера
    document.getElementById('tg-login-button').appendChild(script);
}

export function openProfileModal(userData) {
    // 1. Создаем оверлей
    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay';
    
    // 2. Набиваем HTML-структуру
    overlay.innerHTML = `
    <div class="profile-modal">
        <button class="close-btn" id="modal-close">✕</button>
        <h2>Личный кабинет</h2>
        
        <!-- Сообщение об ошибках/успехе -->
        <div id="modal-message" class="modal-msg" style="display: none;"></div>
        <h3>Личные данные</h3>
        
        <!-- Форма 1: Данные профиля -->
        <form id="profile-info-username-form" class="hidden">
            <div class="input-group">
                <label>Имя пользователя</label>
                <input type="text" name="username" value="${userData.username}" required>
            </div>
            <button type="submit" class="btn-save">Изменить имя пользователя</button>
            <hr>
        </form>
        <form id="profile-info-email-form" class="hidden">
            <div class="input-group">
                <label>Email</label>
                <input type="email" name="email" value="${userData.email}" required>
            </div>
            <button type="submit" class="btn-save">Изменить почту</button>
            <hr>
        </form>

        <form id="profile-tg-link-email" class="hidden">
            <h3>Привязать почту</h3>
            <div class="input-group">
                <label>Email</label>
                <input type="email" name="email" required>
            </div>
            <button type="submit" class="btn-save">Привязать</button>
            <hr>
        </form>

        <form id="profile-tg-link-email-verify2fa" class="hidden">
            <h3>Введите код подтверждения с почты</h3>
            <div class="input-group">
                <label for="2fa-code">Код подтверждения с почты</label>
                <input
                    name="code"
                    type="text"
                    id="2fa-code"
                    placeholder="000000"
                    maxlength="6"
                    inputmode="numeric"
                    pattern="[0-9]{6}"
                    autocomplete="one-time-code"
                >
            </div>
            <button type="submit" class="btn-save">Подтвердить</button>
            <hr>
        </form>

        <form id="profile-tg-link-set-password" class="hidden">
            <h3>Установить пароль</h3>
            <div class="input-group">
                <label>Новый пароль</label>
                <input type="password" name="new_password" required>
            </div>
            <button type="submit" class="btn-save">Установить пароль</button>
            <hr>
        </form>


        <!-- Форма 2: Смена пароля -->
        <form id="profile-password-form" class="hidden">
            <h3>Безопасность</h3>
            <div class="input-group">
                <label>Текущий пароль</label>
                <input type="password" name="old_password" required>
            </div>
            <div class="input-group">
                <label>Новый пароль</label>
                <input type="password" name="new_password" required>
            </div>
            <button type="submit" class="btn-save">Обновить пароль</button>
            <hr>
        </form>


        <!-- Секция Telegram -->
        <h3>Интеграция</h3>
        <div id="tg-widget-container"></div>
    </div>
    `;
    

    document.body.appendChild(overlay);

    const msgDiv = overlay.querySelector('#modal-message');
    const showMsg = (text, isSuccess = false) => {
        msgDiv.innerText = text;
        msgDiv.style.display = 'block';
        msgDiv.className = isSuccess ? 'modal-msg success' : 'modal-msg error';
    };

    // Передаем showMsg в инициализаторы
    initModalEvents(overlay, showMsg, userData);
    initTelegramSection(showMsg);
}

function initModalEvents(modalElement, showMsg, userData) {

    const closeBtn = modalElement.querySelector('#modal-close');
    const emailTGLinkForm = modalElement.querySelector('#profile-tg-link-email');
    const emailVerify2FATGLinkForm = modalElement.querySelector('#profile-tg-link-email-verify2fa');
    const setPasswordTGLinkForm = modalElement.querySelector('#profile-tg-link-set-password');
    const usernameForm = modalElement.querySelector('#profile-info-username-form');
    const emailForm = modalElement.querySelector('#profile-info-email-form');
    const passwordForm = modalElement.querySelector('#profile-password-form');
    let mode = 'verify-email';
    function updateFormVisibility() {
        const hasEmail = Boolean(userData.email && userData.email.trim() !== '');
        const hasPassword = Boolean(userData.has_password);

        // По умолчанию всё скрываем
        emailTGLinkForm.classList.add('hidden');
        emailVerify2FATGLinkForm.classList.add('hidden');
        setPasswordTGLinkForm.classList.add('hidden');
        usernameForm.classList.add('hidden');
        emailForm.classList.add('hidden');
        passwordForm.classList.add('hidden');

        // Логика переключения шагов
        if (!hasEmail && !hasPassword) {
            // Шаг 1: Нет почты -> Запрашиваем почту
            emailTGLinkForm.classList.remove('hidden');
        } else if (hasEmail && !hasPassword) {
            // Шаг 2: Почта есть, нет пароля -> Запрашиваем пароль
            setPasswordTGLinkForm.classList.remove('hidden');
        } else {
            // Шаг 3: Полный аккаунт -> Показываем редактирование и смену пароля
            usernameForm.classList.remove('hidden');
            emailForm.classList.remove('hidden');
            passwordForm.classList.remove('hidden');
        }
    }
    updateFormVisibility();
    // Закрытие модалки
    closeBtn.addEventListener('click', () => modalElement.remove());
    modalElement.addEventListener('click', (e) => {
        if (e.target === modalElement) modalElement.remove();
    });
    emailTGLinkForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const emailVal = emailTGLinkForm.querySelector('input[name="email"]').value;
        
        const response = await setEmail(emailVal);
        if (response.ok) {
            userData.pendingEmail = emailVal;
            mode = 'set-email';
            emailTGLinkForm.classList.add('hidden');
            emailVerify2FATGLinkForm.classList.remove('hidden');
            showMsg('Код отправлен на вашу почту!', true);
        }else if (response.status === 400) {
            const errorData = await response.json().catch(() => ({}));
            showMsg(errorData.error || 'Почта уже занята');
        } else {
            showMsg('Ошибка при отправке почты. Попробуйте позже.');
        }
    });
    emailVerify2FATGLinkForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const codeVal = emailVerify2FATGLinkForm.querySelector('input[name="code"]').value;
        if (mode === 'set-email') {
            const response = await setEmailVerify2FA(userData.pendingEmail, codeVal);
            if (response.ok) {
                userData.email = userData.pendingEmail;
                emailForm.querySelector('input[name="email"]').value = userData.email;
                updateFormVisibility();
                showMsg('Почта успешно подтверждена!', true);
            } else if (response.status === 400) {
                const errorData = await response.json().catch(() => ({}));
                showMsg(errorData.error || 'Неверный или истекший код');
            } else {
                showMsg('Неверный код или срок его действия истек.');
            }
        } else {
            const response = await changeEmailVerify({ email: userData.email, input_code: codeVal, new_email: userData.pendingEmail});
            if (response.ok) {
                userData.email = userData.pendingEmail;
                updateFormVisibility();
                showMsg('Почта успешно изменена!', true);
            } else if (response.status === 400) {
                const errorData = await response.json().catch(() => ({}));
                showMsg(errorData.error || 'Неверный или истекший код');
            } else {
                showMsg('Неверный код или срок его действия истек.');
            }
        }

    });
    setPasswordTGLinkForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const passVal = setPasswordTGLinkForm.querySelector('input[name="new_password"]').value;
        
        const response = await setEmailPassword(userData.email, passVal);
        if (response.ok) {
            userData.has_password = true; // Фиксируем наличие пароля
            updateFormVisibility(); // Переключаем на полный профиль
            showMsg('Пароль успешно установлен!', true);
        } else {
            const errData = await response.json().catch(() => ({}));
            showMsg(errData.error || 'Не удалось установить пароль');
        }
    });
    // Хэндлер для формы личных данных (Имя / Email)
    usernameForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const newUsername = usernameForm.querySelector('input[name="username"]').value;
            
        const response = await changeInfo({ username: newUsername });
        if (response.ok) {
            userData.username = newUsername;
            setupDashboardUI(newUsername);
            showMsg('Данные профиля успешно обновлены!', true);
        } else {
            showMsg('Ошибка при обновлении профиля. Попробуйте позже.');
        }
    });
    emailForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const newEmail = emailForm.querySelector('input[name="email"]').value;
            
        const response = await changeInfo({ email: userData.email, new_email: newEmail });
        if (response.ok) {
            userData.pendingEmail = newEmail;
            emailForm.classList.add('hidden')
            emailVerify2FATGLinkForm.classList.remove('hidden')
        } else if (response.status == 400) {
            showMsg('Почта уже занята');
        }
         else {
            showMsg('Ошибка при обновлении профиля. Попробуйте позже.');
        }
    });
    // Хэндлер для формы смены пароля
    passwordForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const payload = {
            old_password: passwordForm.querySelector('input[name="old_password"]').value,
            new_password: passwordForm.querySelector('input[name="new_password"]').value
        };
            
        const response = await changePassword(payload);
        if (response.ok) {
            showMsg('Пароль успешно изменен!', true);
            passwordForm.reset(); 
        } else {
            const errData = await response.json().catch(() => ({}));
            showMsg(errData.error || 'Не удалось изменить пароль');
        }
    });
}