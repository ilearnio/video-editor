import { stripQuotes } from './videoQuotes'

describe('Services: videoQuotes', () => {
  describe('stripQuotes', () => {
    it('should remove quotation marks', () => {
      expect(stripQuotes('"If you want something, work for it."')).toBe('If you want something, work for it.')
      expect(stripQuotes("'If you want something, work for it.'")).toBe('If you want something, work for it.')
      expect(stripQuotes('“If you want something, work for it.”')).toBe('If you want something, work for it.')
    })

    it('should handle complex cases', () => {
      expect(stripQuotes('"Test text" - John Doe')).toBe('Test text')
      expect(stripQuotes('"Test text" by John Doe')).toBe('Test text')
    })
  })
})
