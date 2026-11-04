/** Pixel-tolerant scroll boundaries also handle fractional layout widths. */
export const categoryRailState = (
    left: number,
    scrollWidth: number,
    clientWidth: number
) => {
    const maximum = Math.max(0, scrollWidth - clientWidth)
    const position = Math.min(maximum, Math.max(0, left))
    return { previous: position > 2, next: maximum - position > 2, maximum }
}
export const categoryRailTarget = (
    key: string,
    left: number,
    scrollWidth: number,
    clientWidth: number
): number | null => {
    const { maximum } = categoryRailState(left, scrollWidth, clientWidth)
    if (key === 'Home') return 0
    if (key === 'End') return maximum
    if (key === 'ArrowLeft') return Math.max(0, left - clientWidth * 0.8)
    if (key === 'ArrowRight') return Math.min(maximum, left + clientWidth * 0.8)
    return null
}
