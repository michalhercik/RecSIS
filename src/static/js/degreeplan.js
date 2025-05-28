setupCheckboxShiftClick();
updateStickyOffset();

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

// Update the sticky offset for the checked-courses-menu based on the header height
function updateStickyOffset() {
    const header = document.querySelector('header');
    const menu = document.getElementById('dp-checked-courses-menu');
    if (header && menu) {
        const height = header.offsetHeight;
        menu.style.top = height + 'px';
    }
}