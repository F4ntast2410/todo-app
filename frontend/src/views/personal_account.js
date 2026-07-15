import { changeInfo, changePassword, checkLinkTg } from "../api/personal_account";

async function initTelegramSection(showMsg) {
    const container = document.getElementById('tg-widget-container');
    if (!container) return;

    try {
        // 1. Запрашиваем у бэкенда текущий статус привязки
        const response = await checkLinkTg();
        if (!response.ok) throw new Error('Не удалось получить статус привязки');
        
        const data = await response.json();

        if (data.Linked) {
            // 2. Если уже привязан — рендерим статус
            renderTelegramLinked(container, data.Username);
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
            <span class="tg-icon">🌀</span>
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
            const response = await fetch('/auth/telegram/link', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(user)
            });

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
    script.setAttribute('data-telegram-login', 'ИМЯ_ТВОЕГО_БОТА'); // Замени на имя своего бота
    script.setAttribute('data-size', 'large');
    script.setAttribute('data-radius', '10');
    script.setAttribute('data-onauth', 'onTelegramAuth(user)'); // Указываем наш глобальный коллбэк
    script.setAttribute('data-request-access', 'write');

    // Вставляем скрипт, чтобы Telegram отрендерил кнопку внутри контейнера
    document.getElementById('tg-login-button').appendChild(script);
}

function openProfileModal(userData) {
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

            <!-- Форма 1: Данные профиля -->
            <form id="profile-info-form">
                <h3>Личные данные</h3>
                <div class="input-group">
                    <label>Имя пользователя</label>
                    <input type="text" name="username" value="${userData.username}" required>
                </div>
                <div class="input-group">
                    <label>Email</label>
                    <input type="email" name="email" value="${userData.email}" required>
                </div>
                <button type="submit" class="btn-save">Сохранить изменения</button>
            </form>

            <hr>

            <!-- Форма 2: Смена пароля -->
            <form id="profile-password-form">
                <h3>Безопасность</h3>
                <div class="input-group">
                    <label>Текущий пароль</label>
                    <input type="password" name="current_password" required>
                </div>
                <div class="input-group">
                    <label>Старый пароль пароль</label>
                    <input type="password" name="old_password" required>
                    <label>Новый пароль</label>
                    <input type="password" name="new_password" required>
                </div>
                <button type="submit" class="btn-save">Обновить пароль</button>
            </form>

            <hr>

            <!-- Секция Telegram -->
            <h3>Интеграция</h3>
            <div id="tg-widget-container"></div>
        </div>
    `;

    document.body.appendChild(overlay);

    // 3. Инициализируем логику и обработчики
    initModalEvents(overlay);
    initTelegramWidget();
}

function initModalEvents(modalElement) {
    const closeBtn = modalElement.querySelector('#modal-close');
    const infoForm = modalElement.querySelector('#profile-info-form');
    const passwordForm = modalElement.querySelector('#profile-password-form');
    const msgDiv = modalElement.querySelector('#modal-message');

    // Функция для вывода уведомлений
    const showMsg = (text, isSuccess = false) => {
        msgDiv.innerText = text;
        msgDiv.style.display = 'block';
        msgDiv.className = isSuccess ? 'modal-msg success' : 'modal-msg error';
    };

    // Закрытие модалки
    closeBtn.addEventListener('click', () => modalElement.remove());
    modalElement.addEventListener('click', (e) => {
        if (e.target === modalElement) modalElement.remove();
    });

    // Хэндлер для формы личных данных (Имя / Email)
    infoForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const payload = {
            username: infoForm.querySelector('input[name="username"]').value,
            email: infoForm.querySelector('input[name="email"]').value
        };
        
        const response = await changeInfo()

        if (response.ok) {
            showMsg('Данные профиля успешно обновлены!', true);
        } else {
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
        
        const response = await changePassword()

        if (response.ok) {
            showMsg('Пароль успешно изменен!', true);
            passwordForm.reset(); 
        } else {
            const errData = await response.json();
            // Перехватываем твою кастомную ошибку "Неверный текущий пароль"
            showMsg(errData.error || 'Не удалось изменить пароль');
        }
    });
}