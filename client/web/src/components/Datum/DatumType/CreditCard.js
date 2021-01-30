import Aux from '../../../hoc/Aux/Aux';
import React from 'react';

export class CreditCard {
  static card(props) {
    return <Aux>
      <h4>{props.datum.desc}</h4>
      <p>Ending with {props.datum.number.slice(
          props.datum.number.length - 4)}</p>
    </Aux>;
  }
}