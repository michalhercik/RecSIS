initAll();

function initAll() {
    initializeTooltips();
    initializePopovers();
    setupCheckboxShiftClick();
    updateStickyOffset();
    updateSmallScreen();
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
    const checkbox = event.target.firstChild.firstChild;
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

// Run on resize
window.addEventListener('resize', updateStickyOffset);

// Update the sticky offset for the checked-courses-menu based on the navbar height
function updateStickyOffset() {
    const navbar = document.getElementById('navbarNav');
    const menu = document.getElementById('checked-courses-menu');
    const height = navbar.offsetHeight;
    menu.style.top = height + 'px';
}


// Update on resize
window.addEventListener('resize', updateSmallScreen);

// Update the smallScreen property in Alpine.js based on window width
function updateSmallScreen() {
    if (window.alpineRef) {
        window.alpineRef.smallScreen = window.innerWidth < 768;
    }
}

