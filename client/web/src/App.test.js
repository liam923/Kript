import React from 'react';
import { render } from '@testing-library/react';
import App from './App';

test('renders datums', () => {
  const { getByText } = render(<App />);
  const headElement = getByText(/My Data/i);
  expect(headElement).toBeInTheDocument();
});
