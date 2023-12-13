import io
import os
from PIL import Image
from typing import Tuple

from langdetect import detect
import pytesseract

from message_types import EvidenceCreatedBody
from services import AShirtRequestsService


def process_content(body: EvidenceCreatedBody):
    ashirt_svc = AShirtRequestsService(
        os.environ.get('ASHIRT_BACKEND_URL', ''),
        os.environ.get('ASHIRT_ACCESS_KEY', ''),
        os.environ.get('ASHIRT_SECRET_KEY', '')
    )
    # gather content
    evidence_content = ashirt_svc.get_evidence_content(
        body.operation_slug, body.evidence_uuid, 'media')
    img = Image.open(io.BytesIO(evidence_content))
    # Run tesseract
    text = pytesseract.image_to_string(
        img, config='--oem 3 --psm 12 -c thresholding_method=2')
    return text