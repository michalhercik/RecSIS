initializeTooltips();

// bootstrap tooltip initialization
function initializeTooltips() {
    const tooltipTriggerList = document.querySelectorAll('[data-bs-toggle="tooltip"]')
    const tooltipList = [...tooltipTriggerList].map(tooltipTriggerEl => new bootstrap.Tooltip(tooltipTriggerEl))
};

// uncheck checkbox by id
function uncheck(id) {
    const checkbox = document.getElementById(id);
    if (checkbox) {
        const clickEvent = new MouseEvent('click', {
            bubbles: true,
            cancelable: true,
        });
        checkbox.dispatchEvent(clickEvent);
    }
}