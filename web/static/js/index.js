// Мобильное меню: открыть/закрыть по кнопке, закрыть по клику на ссылку и по Escape
(() => {
  const header = document.querySelector('.site-header');
  if (!header) return;

  const toggle = header.querySelector('.nav-toggle');
  const nav = header.querySelector('.nav');

  const setOpen = (open) => {
    header.classList.toggle('is-open', open);
    toggle.setAttribute('aria-expanded', String(open));
    toggle.setAttribute('aria-label', open ? 'Закрыть меню' : 'Открыть меню');
  };

  toggle.addEventListener('click', () => {
    setOpen(!header.classList.contains('is-open'));
  });

  nav.addEventListener('click', (e) => {
    if (e.target.closest('a')) setOpen(false);
  });

  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') setOpen(false);
  });
})();
