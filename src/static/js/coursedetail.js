initializeTooltips();

// bootstrap tooltip initialization
function initializeTooltips() {
    // Dispose of existing tooltips
    removeAllTooltips();

    // Initialize new tooltips
    const tooltipTriggerList = document.querySelectorAll('[data-bs-toggle="tooltip"]')
    const tooltipList = [...tooltipTriggerList].map(tooltipTriggerEl => new bootstrap.Tooltip(tooltipTriggerEl))
};

// Dispose of existing tooltips
function removeAllTooltips() {
    var tooltipElements = document.querySelectorAll('.tooltip');
    tooltipElements.forEach(function (tooltipEl) {
        tooltipEl.remove();
    });
};

// uncheck checkbox by id
function uncheck(id) {
    const checkbox = document.getElementById(id);
    console.log('uncheck', id, checkbox);
    if (checkbox) {
        const clickEvent = new MouseEvent('click', {
            bubbles: true,
            cancelable: true,
        });
        checkbox.dispatchEvent(clickEvent);
    }

    const surveyFilters = document.getElementById('survey-filters-form');
    if (surveyFilters) {
        const updateFilters = new CustomEvent('filters-changed', {});
        surveyFilters.dispatchEvent(updateFilters);
    }
}