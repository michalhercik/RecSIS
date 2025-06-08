// Clear expanded- from session storage
function clearExpanded() {
    Object.keys(sessionStorage).forEach(key => {
        if (key.startsWith('expanded-')) {
            sessionStorage.removeItem(key);
        }
    });
}

// Clear the 'cd-tab' session storage item
function clearCDTab() {
    sessionStorage.removeItem('cd-tab');
}

function hideMain() {
    document.getElementById('main-content').style.display = 'none'
}

// Change language to selected
function changeLanguage(lang) {
    let currentUrl = window.location.pathname;
    let queryParams = window.location.search;
    let newUrl = currentUrl.replace(/^\/(cs|en)/, '') || '';
    newUrl = '/' + lang + newUrl + queryParams;
    window.location.href = newUrl;
}

document.addEventListener('htmx:afterSwap', () => {
    // Reinitialize Alpine when HTMX swaps in new content
    console.log('HTMX content swapped, reinitializing Alpine');
    let initAlpine = new Event('alpine:init');
    document.dispatchEvent(initAlpine);
});

window.addEventListener('htmx:historyRestore', () => {
    document.getElementById('main-content').style.display = 'inline';
});