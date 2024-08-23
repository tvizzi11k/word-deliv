document.addEventListener("DOMContentLoaded", () => {
  // Логика для отправки формы
  const form = document.getElementById("contactForm");

  form.addEventListener("submit", (event) => {
    event.preventDefault(); // Останавливаем стандартную отправку формы

    const formData = new FormData(form);

    fetch(form.action, {
      method: form.method,
      body: formData
    })
    .then(response => {
      return response.json().then(data => {
        if (response.ok && data.status === "success") {
          console.log("Форма успешно отправлена, скрываем форму и отображаем сообщение");
          const successMessage = document.getElementById("successMessage");
          const formContainer = document.getElementById("contact-form");

          // Скрываем форму сразу после успешной отправки
          formContainer.style.display = "none";

          // Отображаем сообщение об успешной отправке
          successMessage.style.display = "block";
          setTimeout(() => {
            successMessage.style.display = "none"; // Скрываем сообщение через 3 секунды
            window.location.href = "/"; // Перенаправление на главную страницу
          }, 3000); // Сообщение будет отображаться 3 секунды перед перенаправлением
        } else {
          console.error("Ошибка ответа:", data);
          alert("Ошибка отправки формы: " + (data.message || "Неизвестная ошибка"));
        }
      });
    })
    .catch(error => {
      console.error("Ошибка:", error);
      alert("Ошибка отправки формы.");
    });
  });

  // Логика для открытия и закрытия формы
  const formContainer = document.getElementById("contact-form");
  const openFormButton = document.getElementById("open-form-button");
  const closeFormButton = document.getElementById("close-form");

  openFormButton.addEventListener("click", () => {
    formContainer.style.display = "block";
  });

  closeFormButton.addEventListener("click", () => {
    formContainer.style.display = "none";
  });

  window.addEventListener("click", (event) => {
    if (event.target === formContainer) {
      formContainer.style.display = "none";
    }
  });

  // Логика для бокового меню и бургер-меню
  const burgerMenu = document.getElementById("burger-menu");
  const sidebar = document.getElementById("sidebar");
  const overlay = document.getElementById("sidebar-overlay");
  const burgerMenuInside = document.querySelector('.burger-menu-inside');

  burgerMenu.addEventListener("click", () => {
    burgerMenu.classList.toggle("active");
    sidebar.classList.toggle("active");
    overlay.classList.toggle("active");
  });

  overlay.addEventListener("click", () => {
    burgerMenu.classList.remove("active");
    sidebar.classList.remove("active");
    overlay.classList.remove("active");
  });

  // Закрытие бокового меню при клике на бургер-меню внутри
  burgerMenuInside.addEventListener('click', () => {
    burgerMenu.classList.remove('active');
    sidebar.classList.remove('active');
    overlay.classList.remove('active');
  });

  // Логотип внутри боковой панели
  const sidebarLogo = document.querySelector(".sidebar-logo");
  sidebarLogo.addEventListener("click", () => {
    window.location.href = "/";
  });
});
