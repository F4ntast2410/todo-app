import { changeInfo, changePassword } from "../api/personal_account";

function initTelegramWidget() {
    const container = document.getElementById('tg-widget-container');
    if (!container) return;

    const script = document.createElement('script');
    script.src = "https://telegram.org/js/telegram-widget.js?22";
    script.async = true;
    
    // Переносим настройки из data-атрибутов
    script.setAttribute('data-telegram-login', 'todo_test42_bot');
    script.setAttribute('data-size', 'medium');
    script.setAttribute('data-auth-url', '/auth/telegram/link'); // Куда бэкенд примет callback
    script.setAttribute('data-request-access', 'write');

    // Вставляем скрипт в контейнер — теперь браузер его выполнит и отрендерит кнопку
    container.appendChild(script);
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
        
        const response = await changeInfo

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
        
        const response = await changePassword

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