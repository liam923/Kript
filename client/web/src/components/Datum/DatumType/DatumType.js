import {Password} from './Password';
import {CreditCard} from './CreditCard';

export class DatumType {
  static getType(typeId) {
    switch (typeId) {
      case ("password"):
        return Password;
      case ("creditcards"):
        return CreditCard;
      default:
        return Password;
    }
  }
}