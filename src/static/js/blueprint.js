initAll();

function initAll() {
    initializeTooltips();
    initializePopovers();
    setupCheckboxShiftClick();
}

// bootstrap tooltip initialization
function initializeTooltips() {
    // Dispose of existing tooltips
    var tooltipElements = document.querySelectorAll('.tooltip');
    tooltipElements.forEach(function (tooltipEl) {
        tooltipEl.remove();
    });

    // Initialize new tooltips
    const tooltipTriggerList = document.querySelectorAll('[data-bs-toggle="tooltip"]')
    const tooltipList = [...tooltipTriggerList].map(tooltipTriggerEl => new bootstrap.Tooltip(tooltipTriggerEl))
};

// bootstrap popover initialization
function initializePopovers() {
    // Initialize new popovers
    const popoverTriggerList = document.querySelectorAll('[data-bs-toggle="popover"]')
    const popoverList = [...popoverTriggerList].map(popoverTriggerEl => new bootstrap.Popover(popoverTriggerEl))
};

// shift click for multiple checkboxes
function setupCheckboxShiftClick() {
    let lastChecked = null;
    const checkboxes = document.querySelectorAll('input[type="checkbox"][name="selected"]');
    checkboxes.forEach(checkbox => {
        checkbox.addEventListener('click', function (e) {
            if (lastChecked && lastChecked !== this && e.shiftKey) {
                let inBetween = false;
                checkboxes.forEach(box => {
                    if (box === this || box === lastChecked) {
                        inBetween = !inBetween;
                    }
                    else if (inBetween && box.checked !== this.checked) {
                        box.checked = this.checked;
                        box.dispatchEvent(new Event('change'));
                    }
                });
                if (lastChecked.checked !== this.checked) {
                    lastChecked.checked = this.checked;
                    lastChecked.dispatchEvent(new Event('change'));
                }
            }
            lastChecked = this;
        });
    });
}

// pass through (shift-)click event to checkbox
function handleCircleClick(event) {
    const checkbox = event.target.previousElementSibling;
    if (checkbox) {
        // Create a new MouseEvent, preserving shift key and other properties
        const clickEvent = new MouseEvent('click', {
            bubbles: true,
            cancelable: true,
            shiftKey: event.shiftKey,
        });

        checkbox.dispatchEvent(clickEvent);
    }
}

// update position of course in database using dynamic htmx patch
function sortHxPatch(item, year, semester, position, language) {
    htmx.ajax('PATCH', '/' + language + '/blueprint/course/' + item, {
        target: 'main',
        values: { year: year, semester: semester, position: position + 1 }
    });
}