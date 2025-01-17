((window, document, console) => {

	const loadNextItems = () => {
		const itemNextPlaceholderEls = document.querySelectorAll(".items .item_next_placeholder");

		for (const el of itemNextPlaceholderEls) {
			const distance = (el.offsetTop - document.body.scrollTop);
			console.log("EL", el, distance, document.body.scrollTop, document.body.scrollHeight, el.offsetTop, el.scrollTop);
		}
	};

	const convertFlagImagesToDivEls = () => {
		const flagImageEls = document.querySelectorAll("main .international img.flag");

		if (flagImageEls.length === 0) {
			return;
		}

		for (const el of flagImageEls) {
			const src = el.getAttribute("src");

			if (src) {
				const divEl = document.createElement("div");

				divEl.className = "flag";
				divEl.style.backgroundImage = `url("${ window.encodeURI(src) }")`;

				el.parentElement.replaceChild(divEl, el);
			}
		}
	};

	const getExt = (fileURL, ext) => {
		let extLocal = ext;

		if (!extLocal) {
			const fileURLSegments = fileURL.split(".");

			if (fileURLSegments.length > 0) {
				const extLocal2 = fileURLSegments[fileURLSegments.length - 1];

				if (extLocal2.length >= 2 && extLocal2 <= 4) {
					extLocal = extLocal2;
				}
			}
		}

		if (extLocal.charAt(0) === ".") {
			extLocal = extLocal.slice(1)
		}

		return extLocal
	};

	const loadFile = (fileURL, ext) => {
		fetch(fileURL).then((response) => {
			return response.text();
		}).then((text) => {
			const extLocal = getExt(fileURL, ext);

			switch (extLocal) {
				case "css":
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

	const clearMinHeightForItems = (itemEls) => {
		for (const itemEl of itemEls) {
			itemEl.style.minHeight = "0";
		}
	};

	const setMinHeightForItems = (itemEls) => {
		clearMinHeightForItems(itemEls);

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

	const loadHandler = (() => {
		let inited = false;

		const fileData = [
			{
				u: "https://fonts.googleapis.com/css2?family=Playfair+Display:ital,wght@0,400..900;1,400..900&display=swap",
				e: "css",
			},
			{
				u: "https://fonts.googleapis.com/css2?family=Open+Sans:ital,wght@0,300..800;1,300..800&display=swap",
				e: "css",
			},
		];

		return () => {
			if (!inited) {
				inited = true;
			} else {
				return;
			}

			for (const f of fileData) {
				loadFile(f.u, f.e);
			}

			const itemEls = document.querySelectorAll(".items .item");

			if (itemEls && itemEls.length > 0) {
				window.addEventListener("resize", () => {
					setMinHeightForItems(itemEls);
				});

				document.addEventListener("scroll", () => {
					loadNextItems();
				});

				setMinHeightForItems(itemEls);

				let setMinHeightForItemsIterations = 50;
				const setMinHeightForItemsTimeoutId = window.setInterval(() => {
					setMinHeightForItems(itemEls);

					if (setMinHeightForItemsIterations-- <= 0) {
						window.clearTimeout(setMinHeightForItemsTimeoutId);
					}
				}, 25);
			}

			convertFlagImagesToDivEls();
		};
	})();

	if (document.readyState === "interactive" || document.readyState === "complete") {
		loadHandler();
	}

	window.addEventListener("DOMContentLoaded", () => {
		loadHandler();
	});

	window.addEventListener("load", () => {
		loadHandler();
	});

})(window, document, console);