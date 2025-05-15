// Clear expanded- from session storage
function clearExpanded() {
    Object.keys(sessionStorage).forEach(key => {
        if (key.startsWith('expanded-')) {
            sessionStorage.removeItem(key);
        }
    });
}

// Change language to selected
function changeLanguage(lang) {
    let currentUrl = window.location.pathname;
    let queryParams = window.location.search;
    let newUrl = currentUrl.replace(/^\/(cs|en)/, '') || '';
    newUrl = '/' + lang + newUrl + queryParams;
    window.location.href = newUrl;
}

document.addEventListener('htmx:afterSwap', function (e) {
    // Reinitialize Alpine when HTMX swaps in new content
    if (e.detail.target === document.body) {
        let initAlpine = new Event('alpine:init');
        document.dispatchEvent(initAlpine);
    }
});

window.addEventListener('htmx:historyRestore', () => {
    console.log(document.getElementById('main-content'));
    document.getElementById('main-content').style.display = 'inline';
});