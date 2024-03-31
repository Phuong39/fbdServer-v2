(() => {

	const loadFile = (fileURL, mimeType) => {
		fetch(fileURL).then((response) => {
			return response.text();
		}).then((text) => {
			switch (mimeType) {
				case "text/css":
				{
					const styleEl = document.createElement("style");
					const textNode = document.createTextNode(text);

					styleEl.appendChild(textNode);
					document.head.appendChild(styleEl);
				}
				break;
			}
		}).catch((err) => {
			console.error(err);
		});
	};

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
		const fileData = [
			{
				url: "https://fonts.googleapis.com/css2?family=Open+Sans:ital,wght@0,300..800;1,300..800&display=swap",
				mimeType: "text/css",
			}
		];

		for (const f of fileData) {
			loadFile(f.url, f.mimeType);
		}

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