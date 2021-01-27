import React from 'react';
import { render } from '@testing-library/react';
import App from './App';

test('renders datums', () => {
  const { getByText } = render(<App />);
  const linkElement = getByText(/My Data/i);
  expect(linkElement).toBeInTheDocument();
});
