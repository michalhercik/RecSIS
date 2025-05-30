initAll();

function initAll() {
    initializeTooltips();
    initializePopovers();
    setupCheckboxShiftClick();
    updateStickyOffset();
}

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
    const checkbox = event.target.parentElement.previousElementSibling;
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

// pass through click event to checkbox on small screens
function handleTdClick(event) {
    const checkbox = event.target.querySelector('input[type="checkbox"][name="selected"]');
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
        target: '#blueprint-page',
        swap: 'outerHTML',
        values: { year: year, semester: semester, position: position + 1 }
    });
}

// Update the sticky offset for the checked-courses-menu based on the header height
function updateStickyOffset() {
    const header = document.querySelector('header');
    const menu = document.getElementById('bp-checked-courses-menu');
    if (header && menu) {
        const height = header.offsetHeight;
        menu.style.top = height + 'px';
    }
}