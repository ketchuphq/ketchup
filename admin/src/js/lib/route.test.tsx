import Route from 'lib/route';

describe('Route', function() {
  it('should format', () => {
    let testCases = [
      ['foo!Bar', '/foo-bar'],
      ['foo!!Bar', '/foo-bar'],
      ['/--!--', '/'],
      ['/--!/ab--', '/-/ab'],
    ];
    testCases.map((t) => {
      let actual = Route.format(t[0]);
      expect(actual).toBe(t[1]);
    });
  });

  it('should not format valid', () => {
    let input = '/foo-bar';
    let expected = '/foo-bar';
    let actual = Route.format(input);
    expect(actual).toBe(expected);
  });
});
