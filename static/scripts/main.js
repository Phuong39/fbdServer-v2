(() => {

	const clearMaxHeightForItems = (itemEls) => {
		for (const itemEl of itemEls) {
			delete itemEl.style.minHeight;
		}
	};

	const setMinHeightForItems = (itemEls) => {
		clearMaxHeightForItems(itemEls);

		let minHeight = 0;

		for (const itemEl of itemEls) {
			if (itemEl.clientHeight > minHeight) {
				minHeight = itemEl.clientHeight;
			}
		}

		for (const itemEl of itemEls) {
			itemEl.style.minHeight = minHeight.toString(10) + "px";
		}
	};

	const loadHandler = () => {
		const itemEls = document.querySelectorAll(".items .item");

		if (itemEls && itemEls.length > 0) {
			setMinHeightForItems(itemEls);

			let setMinHeightForItemsIterations = 50;
			const setMinHeightForItemsTimeoutId = window.setInterval(() => {
				setMinHeightForItems(itemEls);

				if (setMinHeightForItemsIterations-- <= 0) {
					window.clearTimeout(setMinHeightForItemsTimeoutId);
				}
			}, 25);
		}
	};

	window.addEventListener("DOMContentLoaded", () => {
		loadHandler();
	});

})();